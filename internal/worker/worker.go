package worker

import (
	"context"
	"fmt"
	"gdtq/internal/broker"
	"sync"
	"time"
)

func StartWorker(ctx context.Context, id int, b *broker.RedisBroker, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d: Ready and polling Redis...\n", id)

	for {
		// Check context at the start of every loop
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Stopping...\n", id)
			return
		default:
			// Pop will now return every 1 second if no task is found
			task, err := b.Pop(ctx)
			if err != nil {
				// This error is expected every 1 second when idle
				continue
			}

			fmt.Printf("Worker %d: Picked up task %d. Finishing...\n", id, task.ID)
			task.Process(id)
			time.Sleep(2 * time.Second)
			fmt.Printf("Worker %d: Completed task %d.\n", id, task.ID)
		}
	}
}
