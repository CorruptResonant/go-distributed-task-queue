package broker

import (
	"context"
	"encoding/json"
	"go-distributed-task-queue/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultQueueName   = "gdtq_tasks"
	DefaultPollTimeout = 1 * time.Second
)

type RedisBroker struct {
	client *redis.Client
	queue  string
}

func NewRedisBroker(addr string) *RedisBroker {
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	return &RedisBroker{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		queue: DefaultQueueName,
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
	results, err := r.client.BRPop(ctx, DefaultPollTimeout, r.queue).Result()
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
