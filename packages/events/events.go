package events

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	EventMarkerCreated = "marker.created"
	EventJobCreated    = "job.created"
	EventJobCompleted  = "job.completed"
	EventJobFailed     = "job.failed"
)

var ErrQueueEmpty = errors.New("queue empty")

type Queue struct {
	client *redis.Client
	key    string
}

func NewQueue(client *redis.Client, key string) *Queue {
	return &Queue{client: client, key: key}
}

func NewRedisClient(addr string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
}

func (q *Queue) Ping(ctx context.Context) error {
	return q.client.Ping(ctx).Err()
}

func (q *Queue) EnqueueJob(ctx context.Context, jobID string) error {
	return q.client.LPush(ctx, q.key, jobID).Err()
}

func (q *Queue) DequeueJob(ctx context.Context, timeout time.Duration) (string, error) {
	result, err := q.client.BRPop(ctx, timeout, q.key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrQueueEmpty
	}
	if err != nil {
		return "", err
	}
	if len(result) != 2 {
		return "", ErrQueueEmpty
	}
	return result[1], nil
}
