package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

func (p *Processor) processAutomationJob(ctx context.Context) error {
	jobIDs, err := p.store.ListPendingAutomationJobIDs(ctx, 1)
	if err != nil {
		return err
	}
	if len(jobIDs) == 0 {
		return nil
	}
	job, claimed, err := p.store.ClaimAutomationJob(ctx, jobIDs[0])
	if err != nil || !claimed {
		return err
	}

	switch job.JobType {
	case stream.AutomationJobTypeYouTubeMetadataSync:
		return p.processYouTubeMetadataJob(ctx, job)
	default:
		if err := p.store.FailAutomationJob(ctx, job.ID, "unsupported automation job type"); err != nil {
			return err
		}
		return errors.New("unsupported automation job type")
	}
}

func (p *Processor) processYouTubeMetadataJob(ctx context.Context, job stream.AutomationJob) error {
	segment, err := p.store.GetPlanSegment(ctx, job.ResourceID)
	if err != nil {
		_ = p.store.FailAutomationJob(ctx, job.ID, err.Error())
		return err
	}
	videoID, canonicalURL, err := normalizeYouTubeURL(segment.YouTubeURL)
	if err != nil {
		_ = p.store.UpdateSegmentMetadata(ctx, segment.ID, "", "", "", "", stream.MetadataStatusFailed, err.Error())
		_ = p.store.FailAutomationJob(ctx, job.ID, err.Error())
		return err
	}

	meta, err := fetchYouTubeMetadata(ctx, canonicalURL)
	if err != nil {
		_ = p.store.UpdateSegmentMetadata(ctx, segment.ID, videoID, "", "", canonicalURL, stream.MetadataStatusFailed, err.Error())
		_ = p.store.FailAutomationJob(ctx, job.ID, err.Error())
		return err
	}

	if err := p.store.UpdateSegmentMetadata(ctx, segment.ID, videoID, meta.Title, meta.AuthorName, canonicalURL, stream.MetadataStatusReady, ""); err != nil {
		_ = p.store.FailAutomationJob(ctx, job.ID, err.Error())
		return err
	}
	return p.store.CompleteAutomationJob(ctx, job.ID)
}

type youtubeOEmbed struct {
	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
}

func fetchYouTubeMetadata(ctx context.Context, canonicalURL string) (youtubeOEmbed, error) {
	endpoint := "https://www.youtube.com/oembed?format=json&url=" + url.QueryEscape(canonicalURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return youtubeOEmbed{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return youtubeOEmbed{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return youtubeOEmbed{}, fmt.Errorf("oembed failed: %s", resp.Status)
	}
	var out youtubeOEmbed
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return youtubeOEmbed{}, err
	}
	return out, nil
}

func normalizeYouTubeURL(raw string) (videoID string, canonicalURL string, err error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", "", err
	}
	switch {
	case strings.Contains(u.Host, "youtu.be"):
		videoID = strings.Trim(u.Path, "/")
	case strings.Contains(u.Host, "youtube.com"):
		videoID = u.Query().Get("v")
		if videoID == "" && strings.HasPrefix(u.Path, "/shorts/") {
			videoID = strings.TrimPrefix(u.Path, "/shorts/")
		}
	}
	if videoID == "" {
		return "", "", errors.New("unsupported YouTube URL")
	}
	return videoID, "https://www.youtube.com/watch?v=" + videoID, nil
}
