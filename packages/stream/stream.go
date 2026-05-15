package stream

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

var ErrNotFound = errors.New("not found")

type Marker struct {
	ID        string    `json:"id"`
	StreamID  string    `json:"stream_id"`
	Label     string    `json:"label"`
	Timestamp string    `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

type Job struct {
	ID            string    `json:"id"`
	MarkerID      string    `json:"marker_id"`
	Status        JobStatus `json:"status"`
	FFmpegCommand string    `json:"ffmpeg_command,omitempty"`
	Error         string    `json:"error,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateMarkerInput struct {
	StreamID  string `json:"stream_id"`
	Label     string `json:"label"`
	Timestamp string `json:"timestamp"`
}

type MarkerJobStore interface {
	CreateMarkerWithJob(ctx context.Context, marker Marker, job Job) error
	GetMarker(ctx context.Context, markerID string) (Marker, error)
	GetJob(ctx context.Context, jobID string) (Job, error)
}

type JobQueue interface {
	EnqueueJob(ctx context.Context, jobID string) error
}

type Service struct {
	store MarkerJobStore
	queue JobQueue
	now   func() time.Time
	idGen func(prefix string) (string, error)
}

func NewService(store MarkerJobStore, queue JobQueue) *Service {
	return &Service{
		store: store,
		queue: queue,
		now:   time.Now().UTC,
		idGen: randomID,
	}
}

func (s *Service) CreateMarker(ctx context.Context, input CreateMarkerInput) (Marker, Job, error) {
	if input.StreamID == "" {
		return Marker{}, Job{}, errors.New("stream_id is required")
	}
	if input.Label == "" {
		return Marker{}, Job{}, errors.New("label is required")
	}
	if input.Timestamp == "" {
		return Marker{}, Job{}, errors.New("timestamp is required")
	}

	markerID, err := s.idGen("marker")
	if err != nil {
		return Marker{}, Job{}, fmt.Errorf("generate marker id: %w", err)
	}

	jobID, err := s.idGen("job")
	if err != nil {
		return Marker{}, Job{}, fmt.Errorf("generate job id: %w", err)
	}

	now := s.now()
	marker := Marker{
		ID:        markerID,
		StreamID:  input.StreamID,
		Label:     input.Label,
		Timestamp: input.Timestamp,
		CreatedAt: now,
	}
	job := Job{
		ID:        jobID,
		MarkerID:  markerID,
		Status:    JobStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.CreateMarkerWithJob(ctx, marker, job); err != nil {
		return Marker{}, Job{}, fmt.Errorf("store marker with job: %w", err)
	}
	if err := s.queue.EnqueueJob(ctx, job.ID); err != nil {
		return Marker{}, Job{}, fmt.Errorf("enqueue job: %w", err)
	}

	return marker, job, nil
}

func randomID(prefix string) (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(buf)), nil
}
