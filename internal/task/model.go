package task

import "time"

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Completed   bool      `json:"completed"`
	SharedWith  []string  `json:"shared_with"`
	CreatedAt   time.Time `json:"created_at"`
	Owner       string    `json:"owner"`
}
