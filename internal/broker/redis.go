package broker

import (
	"context"
	"encoding/json"
	"gdtq/internal/models"
	"time" // Added time import

	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	client *redis.Client
	queue  string
}

func NewRedisBroker(addr string) *RedisBroker {
	return &RedisBroker{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		queue: "gdtq_tasks",
	}
}

func (r *RedisBroker) Push(ctx context.Context, task models.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return r.client.LPush(ctx, r.queue, data).Err()
}

func (r *RedisBroker) Pop(ctx context.Context) (models.Task, error) {
	var task models.Task
	// CHANGE: We use 1*time.Second instead of 0.
	// This ensures the command returns to Go every second to check the context.
	results, err := r.client.BRPop(ctx, 1*time.Second, r.queue).Result()
	if err != nil {
		return task, err
	}
	err = json.Unmarshal([]byte(results[1]), &task)
	return task, err
}

func (r *RedisBroker) GetQueueLength(ctx context.Context) (int64, error) {
	return r.client.LLen(ctx, r.queue).Result()
}

func (r *RedisBroker) Close() error {
	return r.client.Close()
}
