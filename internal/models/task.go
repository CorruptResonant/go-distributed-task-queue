package models

import "fmt"

type Task struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func (t Task) Process(workerID int) {
	fmt.Printf("Worker %d: Executing Task %d [%s] with data: %s\n", workerID, t.ID, t.Type, t.Payload)
}
