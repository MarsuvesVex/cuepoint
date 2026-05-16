package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type runtimeStore interface {
	BuildRuntimeState(ctx context.Context, channelLogin string, now time.Time) (stream.RuntimeState, error)
	GetOperatorProfileByChannelLogin(ctx context.Context, login string) (stream.OperatorProfile, error)
	GetBotSettingsOrDefault(ctx context.Context, userID string) (stream.BotSettings, error)
	GetActivePlanForUser(ctx context.Context, userID string, now time.Time) (stream.StreamPlan, error)
	ListPlanSegments(ctx context.Context, planID string) ([]stream.PlanSegment, error)
	GetOpenStreamSession(ctx context.Context, userID string) (stream.StreamSession, error)
	UpsertStreamSession(ctx context.Context, session stream.StreamSession) error
	UpdateStreamSession(ctx context.Context, session stream.StreamSession) error
	CreateTimelineMarker(ctx context.Context, marker stream.TimelineMarker) error
	FindLatestUnpairedTimelineMarker(ctx context.Context, sessionID, label string) (stream.TimelineMarker, error)
	PairTimelineMarkers(ctx context.Context, startMarkerID, endMarkerID string) error
}

type twitchBridge interface {
	GetChannel(ctx context.Context, userID string) (TwitchChannelState, error)
	GetStream(ctx context.Context, userID string) (TwitchStreamState, error)
	UpdateTitle(ctx context.Context, userID, title string) error
}

type RuntimeService struct {
	store  runtimeStore
	twitch twitchBridge
	now    func() time.Time
	idGen  func(prefix string) (string, error)
}

type TwitchChannelState struct {
	Title string `json:"title"`
}

type TwitchStreamState struct {
	Live      bool      `json:"live"`
	StreamID  string    `json:"streamId"`
	Title     string    `json:"title"`
	StartedAt time.Time `json:"startedAt"`
}

type MarkerRequest struct {
	Label string `json:"label"`
	End   bool   `json:"end"`
}

func NewRuntimeService(store runtimeStore, twitch twitchBridge) *RuntimeService {
	return &RuntimeService{
		store:  store,
		twitch: twitch,
		now:    time.Now().UTC,
		idGen:  streamID,
	}
}

func (s *RuntimeService) GetRuntime(ctx context.Context, channelLogin string) (stream.RuntimeState, error) {
	return s.store.BuildRuntimeState(ctx, channelLogin, s.now())
}

func (s *RuntimeService) SyncSession(ctx context.Context, channelLogin string) (stream.RuntimeState, error) {
	profile, settings, plan, segments, session, hasSession, err := s.loadChannelContext(ctx, channelLogin)
	if err != nil {
		return stream.RuntimeState{}, err
	}

	streamState, err := s.twitch.GetStream(ctx, profile.UserID)
	if err != nil {
		return stream.RuntimeState{}, err
	}

	now := s.now()
	if !streamState.Live {
		if hasSession {
			session.Status = stream.SessionStatusEnded
			endedAt := now
			session.EndedAt = &endedAt
			session.LastSyncedAt = now
			session.UpdatedAt = now
			if session.OriginalTitle != "" {
				_ = s.twitch.UpdateTitle(ctx, profile.UserID, session.OriginalTitle)
			}
			if err := s.store.UpdateStreamSession(ctx, session); err != nil {
				return stream.RuntimeState{}, err
			}
		}
		return s.store.BuildRuntimeState(ctx, channelLogin, now)
	}

	if !hasSession {
		sessionID, err := s.idGen("session")
		if err != nil {
			return stream.RuntimeState{}, err
		}
		session = stream.StreamSession{
			ID:               sessionID,
			UserID:           profile.UserID,
			PlanID:           plan.ID,
			TwitchStreamID:   streamState.StreamID,
			Status:           stream.SessionStatusLive,
			StartedAt:        streamState.StartedAt,
			OriginalTitle:    streamState.Title,
			CurrentTitle:     streamState.Title,
			AutoTitleEnabled: plan.AutoTitleEnabled && settings.AutoTitlesDefaultEnabled,
			LastSyncedAt:     now,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := s.store.UpsertStreamSession(ctx, session); err != nil {
			return stream.RuntimeState{}, err
		}
	} else {
		session.TwitchStreamID = streamState.StreamID
		session.LastSyncedAt = now
		session.UpdatedAt = now
		if session.OriginalTitle == "" {
			session.OriginalTitle = streamState.Title
		}
		if session.CurrentTitle == "" {
			session.CurrentTitle = streamState.Title
		}
		if err := s.store.UpdateStreamSession(ctx, session); err != nil {
			return stream.RuntimeState{}, err
		}
	}

	active := stream.ResolveActiveSegment(plan, session, append([]stream.PlanSegment(nil), segments...), now)
	if active != nil && session.CurrentSegmentID != active.ID {
		session.CurrentSegmentID = active.ID
		started := now
		session.CurrentSegmentStartedAt = &started
		session.UpdatedAt = now
		if err := s.store.UpdateStreamSession(ctx, session); err != nil {
			return stream.RuntimeState{}, err
		}
		if _, err := s.createMarker(ctx, session, *active, stream.MarkerKindSegmentStart, active.SegmentTitle, stream.MarkerSourceSystem); err != nil {
			return stream.RuntimeState{}, err
		}
		if session.AutoTitleEnabled {
			if err := s.applySegmentTitle(ctx, profile.UserID, settings, &session, *active); err != nil {
				return stream.RuntimeState{}, err
			}
		}
	}

	return s.store.BuildRuntimeState(ctx, channelLogin, now)
}

func (s *RuntimeService) ApplyCurrentTitle(ctx context.Context, channelLogin string) (stream.RuntimeState, error) {
	state, err := s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	if err != nil {
		return stream.RuntimeState{}, err
	}
	if !state.IsLive {
		return stream.RuntimeState{}, errors.New("stream is not live")
	}
	if state.ActiveSegment == nil {
		return stream.RuntimeState{}, errors.New("no active segment")
	}
	session := state.Session
	if err := s.applySegmentTitle(ctx, state.Profile.UserID, state.Settings, &session, *state.ActiveSegment); err != nil {
		return stream.RuntimeState{}, err
	}
	return s.store.BuildRuntimeState(ctx, channelLogin, s.now())
}

func (s *RuntimeService) RestoreTitle(ctx context.Context, channelLogin string) (stream.RuntimeState, error) {
	state, err := s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	if err != nil {
		return stream.RuntimeState{}, err
	}
	if !state.IsLive {
		return stream.RuntimeState{}, errors.New("stream is not live")
	}
	if state.Session.OriginalTitle == "" {
		return stream.RuntimeState{}, errors.New("original title is not set")
	}
	if err := s.twitch.UpdateTitle(ctx, state.Profile.UserID, state.Session.OriginalTitle); err != nil {
		return stream.RuntimeState{}, err
	}
	session := state.Session
	session.CurrentTitle = state.Session.OriginalTitle
	session.UpdatedAt = s.now()
	if err := s.store.UpdateStreamSession(ctx, session); err != nil {
		return stream.RuntimeState{}, err
	}
	return s.store.BuildRuntimeState(ctx, channelLogin, s.now())
}

func (s *RuntimeService) ToggleTitles(ctx context.Context, channelLogin string) (stream.RuntimeState, error) {
	state, err := s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	if err != nil {
		return stream.RuntimeState{}, err
	}
	if !state.IsLive {
		return stream.RuntimeState{}, errors.New("stream is not live")
	}
	session := state.Session
	session.AutoTitleEnabled = !session.AutoTitleEnabled
	session.UpdatedAt = s.now()
	if err := s.store.UpdateStreamSession(ctx, session); err != nil {
		return stream.RuntimeState{}, err
	}
	return s.store.BuildRuntimeState(ctx, channelLogin, s.now())
}

func (s *RuntimeService) GetTitleFormat(ctx context.Context, channelLogin string) (string, error) {
	state, err := s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(state.Session.TitleFormatOverride) != "" {
		return state.Session.TitleFormatOverride, nil
	}
	if strings.TrimSpace(state.Settings.DefaultTitleFormat) != "" {
		return state.Settings.DefaultTitleFormat, nil
	}
	return stream.DefaultTitleFormat(), nil
}

func (s *RuntimeService) SetTitleFormat(ctx context.Context, channelLogin, format string) (stream.RuntimeState, error) {
	if err := stream.ValidateTitleFormat(format); err != nil {
		return stream.RuntimeState{}, err
	}
	state, err := s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	if err != nil {
		return stream.RuntimeState{}, err
	}
	if !state.IsLive {
		return stream.RuntimeState{}, errors.New("stream is not live")
	}
	session := state.Session
	session.TitleFormatOverride = format
	session.UpdatedAt = s.now()
	if err := s.store.UpdateStreamSession(ctx, session); err != nil {
		return stream.RuntimeState{}, err
	}
	return s.store.BuildRuntimeState(ctx, channelLogin, s.now())
}

func (s *RuntimeService) ResetTitleFormat(ctx context.Context, channelLogin string) (stream.RuntimeState, error) {
	state, err := s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	if err != nil {
		return stream.RuntimeState{}, err
	}
	if !state.IsLive {
		return stream.RuntimeState{}, errors.New("stream is not live")
	}
	session := state.Session
	session.TitleFormatOverride = ""
	session.UpdatedAt = s.now()
	if err := s.store.UpdateStreamSession(ctx, session); err != nil {
		return stream.RuntimeState{}, err
	}
	return s.store.BuildRuntimeState(ctx, channelLogin, s.now())
}

func (s *RuntimeService) StartSegment(ctx context.Context, channelLogin, segmentID string) (stream.RuntimeState, error) {
	state, err := s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	if err != nil {
		return stream.RuntimeState{}, err
	}
	if !state.IsLive {
		return stream.RuntimeState{}, errors.New("stream is not live")
	}

	var segment *stream.PlanSegment
	for i := range state.Segments {
		if state.Segments[i].ID == segmentID {
			copy := state.Segments[i]
			segment = &copy
			break
		}
	}
	if segment == nil {
		return stream.RuntimeState{}, stream.ErrNotFound
	}

	session := state.Session
	session.CurrentSegmentID = segment.ID
	started := s.now()
	session.CurrentSegmentStartedAt = &started
	session.UpdatedAt = started
	if err := s.store.UpdateStreamSession(ctx, session); err != nil {
		return stream.RuntimeState{}, err
	}
	if _, err := s.createMarker(ctx, session, *segment, stream.MarkerKindSegmentStart, segment.SegmentTitle, stream.MarkerSourceChat); err != nil {
		return stream.RuntimeState{}, err
	}
	if session.AutoTitleEnabled {
		if err := s.applySegmentTitle(ctx, state.Profile.UserID, state.Settings, &session, *segment); err != nil {
			return stream.RuntimeState{}, err
		}
	}
	return s.store.BuildRuntimeState(ctx, channelLogin, s.now())
}

func (s *RuntimeService) AdvanceSegment(ctx context.Context, channelLogin string) (stream.RuntimeState, error) {
	state, err := s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	if err != nil {
		return stream.RuntimeState{}, err
	}
	if !state.IsLive {
		return stream.RuntimeState{}, errors.New("stream is not live")
	}
	next := stream.ResolveNextSegment(state.Plan, state.Session, append([]stream.PlanSegment(nil), state.Segments...), state.Session.CurrentSegmentID)
	if next == nil {
		return stream.RuntimeState{}, errors.New("no next segment")
	}
	if state.ActiveSegment != nil {
		if _, err := s.createMarker(ctx, state.Session, *state.ActiveSegment, stream.MarkerKindSegmentEnd, state.ActiveSegment.SegmentTitle, stream.MarkerSourceChat); err != nil {
			return stream.RuntimeState{}, err
		}
	}
	return s.StartSegment(ctx, channelLogin, next.ID)
}

func (s *RuntimeService) AddMarker(ctx context.Context, channelLogin string, req MarkerRequest) (stream.RuntimeState, error) {
	state, err := s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	if err != nil {
		return stream.RuntimeState{}, err
	}
	if !state.IsLive {
		return stream.RuntimeState{}, errors.New("stream is not live")
	}
	label := strings.TrimSpace(req.Label)
	if label == "" {
		return stream.RuntimeState{}, errors.New("label is required")
	}
	if req.End {
		endMarker, err := s.createMarker(ctx, state.Session, currentSegmentValue(state.ActiveSegment), stream.MarkerKindManualEnd, label, stream.MarkerSourceChat)
		if err != nil {
			return stream.RuntimeState{}, err
		}
		startMarker, err := s.store.FindLatestUnpairedTimelineMarker(ctx, state.Session.ID, label)
		if err == nil && startMarker.ID != endMarker.ID {
			if err := s.store.PairTimelineMarkers(ctx, startMarker.ID, endMarker.ID); err != nil {
				return stream.RuntimeState{}, err
			}
		}
		return s.store.BuildRuntimeState(ctx, channelLogin, s.now())
	}
	if _, err := s.createMarker(ctx, state.Session, currentSegmentValue(state.ActiveSegment), stream.MarkerKindManualStart, label, stream.MarkerSourceChat); err != nil {
		return stream.RuntimeState{}, err
	}
	return s.store.BuildRuntimeState(ctx, channelLogin, s.now())
}

func (s *RuntimeService) loadChannelContext(ctx context.Context, channelLogin string) (stream.OperatorProfile, stream.BotSettings, stream.StreamPlan, []stream.PlanSegment, stream.StreamSession, bool, error) {
	profile, err := s.store.GetOperatorProfileByChannelLogin(ctx, channelLogin)
	if err != nil {
		return stream.OperatorProfile{}, stream.BotSettings{}, stream.StreamPlan{}, nil, stream.StreamSession{}, false, err
	}
	settings, err := s.store.GetBotSettingsOrDefault(ctx, profile.UserID)
	if err != nil {
		return stream.OperatorProfile{}, stream.BotSettings{}, stream.StreamPlan{}, nil, stream.StreamSession{}, false, err
	}
	plan, err := s.store.GetActivePlanForUser(ctx, profile.UserID, s.now())
	if err != nil {
		return stream.OperatorProfile{}, stream.BotSettings{}, stream.StreamPlan{}, nil, stream.StreamSession{}, false, err
	}
	segments, err := s.store.ListPlanSegments(ctx, plan.ID)
	if err != nil {
		return stream.OperatorProfile{}, stream.BotSettings{}, stream.StreamPlan{}, nil, stream.StreamSession{}, false, err
	}
	session, err := s.store.GetOpenStreamSession(ctx, profile.UserID)
	if err != nil {
		if errors.Is(err, stream.ErrNotFound) {
			return profile, settings, plan, segments, stream.StreamSession{}, false, nil
		}
		return stream.OperatorProfile{}, stream.BotSettings{}, stream.StreamPlan{}, nil, stream.StreamSession{}, false, err
	}
	return profile, settings, plan, segments, session, true, nil
}

func (s *RuntimeService) applySegmentTitle(ctx context.Context, userID string, settings stream.BotSettings, session *stream.StreamSession, segment stream.PlanSegment) error {
	format := settings.DefaultTitleFormat
	if strings.TrimSpace(session.TitleFormatOverride) != "" {
		format = session.TitleFormatOverride
	}
	title := stream.BuildSegmentTitle(format, session.OriginalTitle, stream.SegmentDisplayTitle(segment))
	if err := s.twitch.UpdateTitle(ctx, userID, title); err != nil {
		return err
	}
	session.CurrentTitle = title
	session.UpdatedAt = s.now()
	return s.store.UpdateStreamSession(ctx, *session)
}

func (s *RuntimeService) createMarker(ctx context.Context, session stream.StreamSession, segment stream.PlanSegment, kind stream.MarkerKind, label string, source stream.MarkerSource) (stream.TimelineMarker, error) {
	id, err := s.idGen("timeline")
	if err != nil {
		return stream.TimelineMarker{}, err
	}
	now := s.now()
	marker := stream.TimelineMarker{
		ID:                   id,
		UserID:               session.UserID,
		PlanID:               session.PlanID,
		SessionID:            session.ID,
		SegmentID:            segment.ID,
		Kind:                 kind,
		Label:                label,
		StreamElapsedSeconds: stream.StreamElapsedSeconds(session, now),
		Source:               source,
		CreatedAt:            now,
	}
	if err := s.store.CreateTimelineMarker(ctx, marker); err != nil {
		return stream.TimelineMarker{}, err
	}
	return marker, nil
}

func currentSegmentValue(segment *stream.PlanSegment) stream.PlanSegment {
	if segment == nil {
		return stream.PlanSegment{}
	}
	return *segment
}

func streamID(prefix string) (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(buf)), nil
}

type TwitchBridgeClient struct {
	baseURL      string
	headerName   string
	serviceToken string
	client       *http.Client
}

func NewTwitchBridgeClient(baseURL, headerName, serviceToken string, client *http.Client) *TwitchBridgeClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &TwitchBridgeClient{
		baseURL:      strings.TrimRight(baseURL, "/"),
		headerName:   headerName,
		serviceToken: serviceToken,
		client:       client,
	}
}

func (c *TwitchBridgeClient) GetChannel(ctx context.Context, userID string) (TwitchChannelState, error) {
	var out TwitchChannelState
	if err := c.getJSON(ctx, fmt.Sprintf("%s/api/internal/twitch/users/%s/channel", c.baseURL, userID), &out); err != nil {
		return TwitchChannelState{}, err
	}
	return out, nil
}

func (c *TwitchBridgeClient) GetStream(ctx context.Context, userID string) (TwitchStreamState, error) {
	var out TwitchStreamState
	if err := c.getJSON(ctx, fmt.Sprintf("%s/api/internal/twitch/users/%s/stream", c.baseURL, userID), &out); err != nil {
		return TwitchStreamState{}, err
	}
	return out, nil
}

func (c *TwitchBridgeClient) UpdateTitle(ctx context.Context, userID, title string) error {
	body, err := json.Marshal(map[string]string{"title": title})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/internal/twitch/users/%s/title", c.baseURL, userID), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyAuth(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("update title failed: %s", resp.Status)
	}
	return nil
}

func (c *TwitchBridgeClient) getJSON(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	c.applyAuth(req)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request %s failed: %s", url, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *TwitchBridgeClient) applyAuth(req *http.Request) {
	if c.headerName != "" && c.serviceToken != "" {
		req.Header.Set(c.headerName, c.serviceToken)
	}
}
