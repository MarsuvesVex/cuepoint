package stream

import (
	"context"
	"errors"
	"testing"
)

type testStore struct {
	marker Marker
	job    Job
	err    error
}

func (s *testStore) CreateMarkerWithJob(_ context.Context, marker Marker, job Job) error {
	if s.err != nil {
		return s.err
	}
	s.marker = marker
	s.job = job
	return nil
}

func (s *testStore) GetMarker(context.Context, string) (Marker, error) { return Marker{}, nil }
func (s *testStore) GetJob(context.Context, string) (Job, error)       { return Job{}, nil }

type testQueue struct {
	jobID string
	err   error
}

func (q *testQueue) EnqueueJob(_ context.Context, jobID string) error {
	if q.err != nil {
		return q.err
	}
	q.jobID = jobID
	return nil
}

func TestCreateMarker(t *testing.T) {
	store := &testStore{}
	queue := &testQueue{}
	service := NewService(store, queue)
	service.idGen = func(prefix string) (string, error) { return prefix + "-id", nil }

	marker, job, err := service.CreateMarker(context.Background(), CreateMarkerInput{
		StreamID:  "stream",
		Label:     "clip",
		Timestamp: "00:01:00",
	})
	if err != nil {
		t.Fatalf("CreateMarker returned error: %v", err)
	}

	if marker.ID != "marker-id" {
		t.Fatalf("marker ID = %q", marker.ID)
	}
	if job.ID != "job-id" {
		t.Fatalf("job ID = %q", job.ID)
	}
	if queue.jobID != "job-id" {
		t.Fatalf("queue jobID = %q", queue.jobID)
	}
}

func TestCreateMarkerQueueError(t *testing.T) {
	store := &testStore{}
	queue := &testQueue{err: errors.New("boom")}
	service := NewService(store, queue)
	service.idGen = func(prefix string) (string, error) { return prefix + "-id", nil }

	_, _, err := service.CreateMarker(context.Background(), CreateMarkerInput{
		StreamID:  "stream",
		Label:     "clip",
		Timestamp: "00:01:00",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
