package worker

import (
	"encoding/json"
	"log"

	"github.com/Greeg-g/todo-api/internal/cache"
	"github.com/Greeg-g/todo-api/internal/email"
	"github.com/Greeg-g/todo-api/internal/model"
)

const (
	emailPrimary   = "greg.arc03@gmail.com"
	emailSecondary = "gardc.arc03@gmail.com"
)

// Listens for deadline alerts and sends email notifications
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

		subject := "🔔 Upcoming task deadline reminder 🔔"

		// plain text fallback
		var when string
		if task.Deadline != nil {
			when = task.Deadline.Format("02/01/2006 15:04")
		} else {
			when = "(no deadline specified)"
		}
		textBody := "Hello " + task.Owner + ",\n\n" +
			"This is a friendly reminder that the task '" + task.Title + "' is due at " + when + ".\n\n" +
			"Please mark it as completed once done.\n\n" +
			"— The To-Do App"

		// HTML body (simple, inline styles)
		htmlBody := `
		<!doctype html>
			<html>
				<head>
					<meta charset="utf-8">
					<meta name="viewport" content="width=device-width,initial-scale=1">
					<style>
						.card { max-width:600px; margin:20px auto; padding:20px; font-family:Arial, sans-serif; border-radius:8px; background:#f9f9fb; }
						.title { font-size:18px; font-weight:700; color:#333 }
						.meta { color:#666; margin-top:8px }
						.cta { display:inline-block; margin-top:16px; padding:10px 16px; background:#2563eb; color:#fff; text-decoration:none; border-radius:6px }
					</style>
				</head>
				<body>
					<div class="card">
						<div class="title">Task deadline reminder</div>
						<p>Hi <strong>` + task.Owner + `</strong>,</p>
						<p>Your task: <strong>` + task.Title + `</strong></p>
						<p class="meta">Due at: <strong>` + when + `</strong></p>
						<p>If you've already completed it, you can safely ignore this message.</p>
						<a class="cta" href="#">Open your tasks</a>
						<p style="color:#999; font-size:12px; margin-top:12px">Sent by To-Do App</p>
					</div>
				</body>
			</html>
		`

		to := mapUserToEmail(task.Owner)
		if err := email.SendHTML(to, subject, textBody, htmlBody); err != nil {
			log.Println("Error sending email:", err)
		} else {
			log.Printf("Alert email sent to %s for task ID %d\n", to, task.ID)
		}
	}
}

// Maps a username to an email address
func mapUserToEmail(username string) string {
	switch username {
	case "greg":
		return emailPrimary
	default:
		return emailSecondary
	}
}
