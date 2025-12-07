package model

import (
	"time"

	"github.com/lib/pq"
)

type Task struct {
	ID          int64          `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Completed   bool           `json:"completed"`
	SharedWith  pq.StringArray `json:"shared_with" gorm:"type:text[]"`
	CreatedAt   time.Time      `json:"created_at"`
	Deadline    *time.Time     `json:"deadline" gorm:"type:timestamptz;default:NULL"`
	Owner       string         `json:"owner"`
}
