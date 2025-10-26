package main

import (
	"github.com/Greeg-g/todo-api/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.SetupRoutes(r)

	r.Run(":8080")
}
