package api

import (
	"encoding/json"
	"go-distributed-task-queue/internal/broker"
	"go-distributed-task-queue/internal/models"
	"net/http"
)

type TaskHandler struct {
	Broker *broker.RedisBroker
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Broker.Push(r.Context(), t); err != nil {
		http.Error(w, "Failed to enqueue task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task accepted"})
}

func StatsHandler(b *broker.RedisBroker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		length, err := b.GetQueueLength(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"queue_depth": length,
			"status":      "active",
		})
	}
}
