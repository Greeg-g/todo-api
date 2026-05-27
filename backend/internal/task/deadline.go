package task

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Greeg-g/todo-api/internal/cache"
	"github.com/Greeg-g/todo-api/internal/database"
	"github.com/Greeg-g/todo-api/internal/model"
)

// Enqueues tasks with upcoming deadlines into Redis for alerts
func EnqueueUpcomingDeadlines() {
	var tasks []model.Task
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	now := time.Now().In(loc)
	soon := now.Add(time.Hour)

	if err := database.DB.Where("deadline BETWEEN ? AND ? AND completed = false", now, soon).Find(&tasks).Error; err != nil {
		log.Println("Error retrieving upcoming deadlines:", err)
		return
	}

	for _, task := range tasks {
		payload, _ := json.Marshal(task)
		cache.RDB.LPush(cache.Ctx, "deadline_alerts", payload)
	}
}
