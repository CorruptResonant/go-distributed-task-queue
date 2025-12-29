package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go-distributed-task-queue/internal/api"
	"go-distributed-task-queue/internal/broker"
	"go-distributed-task-queue/internal/worker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379" // Local fallback
	}

	redisBroker := broker.NewRedisBroker(redisAddr)

	var wg sync.WaitGroup
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker.StartWorker(ctx, i, redisBroker, &wg)
	}

	mux := http.NewServeMux()
	mux.Handle("/enqueue", &api.TaskHandler{Broker: redisBroker})
	mux.HandleFunc("/stats", api.StatsHandler(redisBroker))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Println("GDTQ Backend running on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server error: %s\n", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\nShutdown signal received. Waiting for in-flight tasks...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(shutdownCtx)
	wg.Wait()
	redisBroker.Close()

	fmt.Println("GDTQ exited cleanly.")
}
