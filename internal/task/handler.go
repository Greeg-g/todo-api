package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Greeg-g/todo-api/internal/cache"
	"github.com/Greeg-g/todo-api/internal/database"
	"github.com/Greeg-g/todo-api/internal/model"
	"github.com/gin-gonic/gin"
)

var tasks = []model.Task{}

func RegisterRoutes(r *gin.Engine) {
	taskGroup := r.Group("/tasks")
	{
		taskGroup.GET("/", getAllTasks)
		taskGroup.POST("/create", createTask)
		taskGroup.POST("/complete/:id", completeTask)
		taskGroup.DELETE("/delete/:id", deleteTask)
		taskGroup.GET("/category/:category", getCategoryTasks)
		taskGroup.POST("/share/:id", shareTask)
		taskGroup.GET("/shared/:user", getSharedTasks)
		taskGroup.GET("/recent", getRecentTasks)
		taskGroup.GET("/owner/:owner", getTasksByOwner)
		taskGroup.GET("check-deadlines", checkDeadlines)
	}
}

func getAllTasks(c *gin.Context) {
	var tasks []model.Task
	if err := database.DB.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var newTask model.Task
	err := c.ShouldBindJSON(&newTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask.CreatedAt = time.Now()
	if err := database.DB.Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, newTask)
}

func completeTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var task model.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.Completed = true
	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := database.DB.Delete(&model.Task{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func getCategoryTasks(c *gin.Context) {
	category := c.Param("category")

	if err := database.DB.Where("LOWER(category) = LOWER(?)", category).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func shareTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		User string `json:"user" binding:"required"`
	}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task model.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if task.Owner == req.User {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Owner cannot be added to shared_with"})
		return
	}

	for _, user := range task.SharedWith {
		if strings.EqualFold(user, req.User) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already has access to this task"})
			return
		}
	}

	task.SharedWith = append(task.SharedWith, req.User)
	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func getSharedTasks(c *gin.Context) {
	user := c.Param("user")
	var tasks []model.Task

	err := database.DB.Where("? = ANY (shared_with)", user).Find(&tasks).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks shared with this user"})
	}

	c.JSON(http.StatusOK, tasks)
}

func getTasksByOwner(c *gin.Context) {
	owner := c.Param("owner")
	var tasks []model.Task
	if err := database.DB.Where("LOWER(owner) = LOWER(?)", owner).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func getRecentTasks(c *gin.Context) {
	val, err := cache.RDB.Get(cache.Ctx, "recent_tasks").Result()
	if err == nil {
		c.Data(http.StatusOK, "application/json", []byte(val))
		return
	}

	var tasks []model.Task
	if err := database.DB.Order("created_at desc").Limit(10).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	jsonData, _ := json.Marshal(tasks)
	cache.RDB.Set(cache.Ctx, "recent_tasks", jsonData, time.Minute*10)

	c.JSON(http.StatusOK, tasks)
}

func checkDeadlines(c *gin.Context) {
	go EnqueueUpcomingDeadlines()
	c.JSON(http.StatusOK, gin.H{"message": "Deadline check initiated"})
}
