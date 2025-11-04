package worker

import (
	"encoding/json"
	"log"

	"github.com/Greeg-g/todo-api/internal/cache"
	"github.com/Greeg-g/todo-api/internal/email"
	"github.com/Greeg-g/todo-api/internal/model"
)

func StartDeadlineWorker() {
	log.Println("Deadline worker started, waiting for tasks...")

	for {
		result, err := cache.RDB.BRPop(cache.Ctx, 0, "deadline_alerts").Result()
		if err != nil {
			log.Println("Error retrieving task from Redis:", err)
			continue
		}

		var task model.Task
		if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
			log.Println("Error unmarshaling task:", err)
			continue
		}

		subject := "🔔 Sua tarefa está perto do prazo!! 🔔"
		body := "Olá " + task.Owner + ",\n\nA tarefa '" + task.Title + "' vence em " +
			task.Deadline.Format("02/01/2006 15:04") + ".\n\nNão se esqueça de concluí-la!"

		to := mapUserToEmail(task.Owner)
		if err := email.Send(to, subject, body); err != nil {
			log.Println("Error sending email:", err)
		} else {
			log.Printf("Alert email sent to %s for task ID %d\n", to, task.ID)
		}
	}
}

func mapUserToEmail(username string) string {
	switch username {
	case "greg":
		return "greg.arc03@gmail.com"
	default:
		return "gardc.arc03@gmail.com"
	}
}
