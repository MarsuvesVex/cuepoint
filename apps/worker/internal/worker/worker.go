package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MarsuvesVex/cuepoint/packages/events"
	"github.com/MarsuvesVex/cuepoint/packages/ffmpeg"
	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type Store interface {
	ListPendingJobIDs(ctx context.Context, limit int) ([]string, error)
	ClaimJob(ctx context.Context, jobID string) (stream.Job, stream.Marker, bool, error)
	CompleteJob(ctx context.Context, jobID, command string) error
	FailJob(ctx context.Context, jobID, reason string) error
}

type Queue interface {
	EnqueueJob(ctx context.Context, jobID string) error
	DequeueJob(ctx context.Context, timeout time.Duration) (string, error)
}

type Processor struct {
	store        Store
	queue        Queue
	blockTimeout time.Duration
	logger       *log.Logger
}

func NewProcessor(store Store, queue Queue, blockTimeout time.Duration, logger *log.Logger) *Processor {
	if logger == nil {
		logger = log.Default()
	}
	return &Processor{
		store:        store,
		queue:        queue,
		blockTimeout: blockTimeout,
		logger:       logger,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	if err := p.requeuePending(ctx); err != nil {
		return fmt.Errorf("requeue pending: %w", err)
	}

	for {
		if ctx.Err() != nil {
			return nil
		}

		jobID, err := p.queue.DequeueJob(ctx, p.blockTimeout)
		if err != nil {
			if errors.Is(err, events.ErrQueueEmpty) || errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("dequeue job: %w", err)
		}

		if err := p.processJob(ctx, jobID); err != nil {
			p.logger.Printf("process job %s: %v", jobID, err)
		}
	}
}

func (p *Processor) requeuePending(ctx context.Context) error {
	jobIDs, err := p.store.ListPendingJobIDs(ctx, 100)
	if err != nil {
		return err
	}
	for _, jobID := range jobIDs {
		if err := p.queue.EnqueueJob(ctx, jobID); err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) processJob(ctx context.Context, jobID string) error {
	job, marker, claimed, err := p.store.ClaimJob(ctx, jobID)
	if err != nil {
		return err
	}
	if !claimed {
		return nil
	}

	command := ffmpeg.BuildSegmentCommand(marker, job)
	if command == "" {
		failErr := p.store.FailJob(ctx, jobID, "empty ffmpeg command")
		if failErr != nil {
			return fmt.Errorf("store failure after empty command: %w", failErr)
		}
		return errors.New("empty ffmpeg command")
	}

	return p.store.CompleteJob(ctx, jobID, command)
}
