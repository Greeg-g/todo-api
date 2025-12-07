package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"

	"github.com/Greeg-g/todo-api/internal/cache"
	"github.com/Greeg-g/todo-api/internal/database"
	"github.com/Greeg-g/todo-api/internal/routes"
	"github.com/Greeg-g/todo-api/internal/task"
	"github.com/Greeg-g/todo-api/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

// Main entry point for the API server
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	log.Println("SMTP_USER:", os.Getenv("SMTP_USER"))
	log.Println("SMTP_PASS:", os.Getenv("SMTP_PASS"))

	database.Connect()
	cache.Connect()

	go worker.StartDeadlineWorker()

	startScheduler()

	r := gin.Default()

	// Permitir CORS para o frontend em localhost:3000
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))
	routes.SetupRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Starts the cron scheduler for periodic deadline checks
func startScheduler() {
	c := cron.New()
	c.AddFunc("@every 5m", func() {
		log.Println("Checking for upcoming task deadlines...")
		task.EnqueueUpcomingDeadlines()
	})
	c.Start()
}
