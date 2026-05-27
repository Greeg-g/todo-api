package routes

import (
	"github.com/Greeg-g/todo-api/internal/auth"
	"github.com/Greeg-g/todo-api/internal/cache"
	"github.com/Greeg-g/todo-api/internal/task"
	"github.com/gin-gonic/gin"
)

// Sets up all API routes for the application
func SetupRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth.RegisterRoutes(r)
	task.RegisterRoutes(r)
	cache.RegisterRoutes(r)
}
