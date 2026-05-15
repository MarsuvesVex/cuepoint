package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}

	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}

func (s *Store) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *Store) Bootstrap(ctx context.Context) error {
	const schema = `
CREATE TABLE IF NOT EXISTS markers (
    id TEXT PRIMARY KEY,
    stream_id TEXT NOT NULL,
    label TEXT NOT NULL,
    timestamp TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    marker_id TEXT NOT NULL REFERENCES markers(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    ffmpeg_command TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
`
	_, err := s.pool.Exec(ctx, schema)
	return err
}

func (s *Store) CreateMarkerWithJob(ctx context.Context, marker stream.Marker, job stream.Job) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `
INSERT INTO markers (id, stream_id, label, timestamp, created_at)
VALUES ($1, $2, $3, $4, $5)
`, marker.ID, marker.StreamID, marker.Label, marker.Timestamp, marker.CreatedAt); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
INSERT INTO jobs (id, marker_id, status, ffmpeg_command, error, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`, job.ID, job.MarkerID, string(job.Status), job.FFmpegCommand, job.Error, job.CreatedAt, job.UpdatedAt); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) GetMarker(ctx context.Context, markerID string) (stream.Marker, error) {
	row := s.pool.QueryRow(ctx, `
SELECT id, stream_id, label, timestamp, created_at
FROM markers
WHERE id = $1
`, markerID)

	var marker stream.Marker
	err := row.Scan(&marker.ID, &marker.StreamID, &marker.Label, &marker.Timestamp, &marker.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return stream.Marker{}, stream.ErrNotFound
	}
	return marker, err
}

func (s *Store) GetJob(ctx context.Context, jobID string) (stream.Job, error) {
	row := s.pool.QueryRow(ctx, `
SELECT id, marker_id, status, ffmpeg_command, error, created_at, updated_at
FROM jobs
WHERE id = $1
`, jobID)

	var job stream.Job
	var status string
	err := row.Scan(&job.ID, &job.MarkerID, &status, &job.FFmpegCommand, &job.Error, &job.CreatedAt, &job.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return stream.Job{}, stream.ErrNotFound
	}
	job.Status = stream.JobStatus(status)
	return job, err
}

func (s *Store) ListPendingJobIDs(ctx context.Context, limit int) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id
FROM jobs
WHERE status = $1
ORDER BY created_at ASC
LIMIT $2
`, string(stream.JobStatusPending), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (s *Store) ClaimJob(ctx context.Context, jobID string) (stream.Job, stream.Marker, bool, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return stream.Job{}, stream.Marker{}, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var job stream.Job
	var marker stream.Marker
	var status string
	var updatedAt time.Time
	err = tx.QueryRow(ctx, `
UPDATE jobs
SET status = $2, updated_at = NOW()
WHERE id = $1 AND status = $3
RETURNING id, marker_id, status, ffmpeg_command, error, created_at, updated_at
`, jobID, string(stream.JobStatusRunning), string(stream.JobStatusPending)).Scan(
		&job.ID,
		&job.MarkerID,
		&status,
		&job.FFmpegCommand,
		&job.Error,
		&job.CreatedAt,
		&updatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return stream.Job{}, stream.Marker{}, false, nil
	}
	if err != nil {
		return stream.Job{}, stream.Marker{}, false, err
	}
	job.Status = stream.JobStatus(status)
	job.UpdatedAt = updatedAt

	err = tx.QueryRow(ctx, `
SELECT id, stream_id, label, timestamp, created_at
FROM markers
WHERE id = $1
`, job.MarkerID).Scan(&marker.ID, &marker.StreamID, &marker.Label, &marker.Timestamp, &marker.CreatedAt)
	if err != nil {
		return stream.Job{}, stream.Marker{}, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return stream.Job{}, stream.Marker{}, false, err
	}

	return job, marker, true, nil
}

func (s *Store) CompleteJob(ctx context.Context, jobID, command string) error {
	tag, err := s.pool.Exec(ctx, `
UPDATE jobs
SET status = $2, ffmpeg_command = $3, error = '', updated_at = NOW()
WHERE id = $1
`, jobID, string(stream.JobStatusCompleted), command)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return stream.ErrNotFound
	}
	return nil
}

func (s *Store) FailJob(ctx context.Context, jobID, reason string) error {
	tag, err := s.pool.Exec(ctx, `
UPDATE jobs
SET status = $2, error = $3, updated_at = NOW()
WHERE id = $1
`, jobID, string(stream.JobStatusFailed), reason)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return stream.ErrNotFound
	}
	return nil
}
