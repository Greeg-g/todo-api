package worker

import (
	"encoding/json"
	"log"

	"github.com/Greeg-g/todo-api/internal/cache"
	"github.com/Greeg-g/todo-api/internal/email"
	"github.com/Greeg-g/todo-api/internal/model"
)

func StartDeadlineWorker() {
	log.Println("⏳ Worker de prazos iniciado")

	for {
		result, err := cache.RDB.BRPop(cache.Ctx, 0, "deadline_alerts").Result()
		if err != nil {
			log.Println("Erro ao buscar da fila:", err)
			continue
		}

		var task model.Task
		if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
			log.Println("Erro ao decodificar tarefa:", err)
			continue
		}

		subject := "🔔 Sua tarefa está perto do prazo!"
		body := "Olá " + task.Owner + ",\n\nA tarefa '" + task.Title + "' vence em " +
			task.Deadline.Format("02/01/2006 15:04") + ".\n\nNão se esqueça de concluí-la!"

		to := mapUserToEmail(task.Owner)
		if err := email.Send(to, subject, body); err != nil {
			log.Println("Erro ao enviar e-mail:", err)
		} else {
			log.Printf("✅ E-mail enviado para %s sobre tarefa '%s'", to, task.Title)
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
