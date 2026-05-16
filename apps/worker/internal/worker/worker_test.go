package worker

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"
	"time"

	"github.com/MarsuvesVex/cuepoint/packages/events"
	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type fakeStore struct {
	pending []string
	job     stream.Job
	marker  stream.Marker
	done    string
}

func (s *fakeStore) ListPendingJobIDs(context.Context, int) ([]string, error) {
	return s.pending, nil
}

func (s *fakeStore) ClaimJob(context.Context, string) (stream.Job, stream.Marker, bool, error) {
	return s.job, s.marker, true, nil
}

func (s *fakeStore) CompleteJob(_ context.Context, _ string, command string) error {
	s.done = command
	return nil
}

func (s *fakeStore) FailJob(context.Context, string, string) error { return nil }
func (s *fakeStore) ListPendingAutomationJobIDs(context.Context, int) ([]string, error) {
	return nil, nil
}
func (s *fakeStore) ClaimAutomationJob(context.Context, string) (stream.AutomationJob, bool, error) {
	return stream.AutomationJob{}, false, nil
}
func (s *fakeStore) GetPlanSegment(context.Context, string) (stream.PlanSegment, error) {
	return stream.PlanSegment{}, nil
}
func (s *fakeStore) UpdateSegmentMetadata(context.Context, string, string, string, string, string, stream.MetadataStatus, string) error {
	return nil
}
func (s *fakeStore) CompleteAutomationJob(context.Context, string) error { return nil }
func (s *fakeStore) FailAutomationJob(context.Context, string, string) error {
	return nil
}

type fakeQueue struct {
	dequeued []string
	enqueued []string
}

func (q *fakeQueue) EnqueueJob(_ context.Context, jobID string) error {
	q.enqueued = append(q.enqueued, jobID)
	return nil
}

func (q *fakeQueue) DequeueJob(ctx context.Context, _ time.Duration) (string, error) {
	if len(q.dequeued) == 0 {
		<-ctx.Done()
		return "", events.ErrQueueEmpty
	}
	jobID := q.dequeued[0]
	q.dequeued = q.dequeued[1:]
	return jobID, nil
}

func TestProcessJob(t *testing.T) {
	store := &fakeStore{
		job:    stream.Job{ID: "job-1"},
		marker: stream.Marker{ID: "marker-1", StreamID: "stream", Timestamp: "00:00:10"},
	}
	queue := &fakeQueue{dequeued: []string{"job-1"}}
	processor := NewProcessor(store, queue, time.Millisecond, log.New(io.Discard, "", 0))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	err := processor.Run(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Fatalf("Run returned error: %v", err)
	}
	if store.done == "" {
		t.Fatal("expected ffmpeg command to be stored")
	}
}
