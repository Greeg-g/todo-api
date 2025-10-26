package task

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var tasks = []Task{}

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
	}
}

func getAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var newTask Task
	err := c.ShouldBindJSON(&newTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask.ID = getNextID()
	newTask.CreatedAt = time.Now()
	tasks = append(tasks, newTask)
	c.JSON(http.StatusCreated, newTask)
}

func completeTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			if task.Completed {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Task already completed"})
				return
			}
			tasks[i].Completed = true
			c.JSON(http.StatusOK, tasks[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func deleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func getCategoryTasks(c *gin.Context) {
	category := c.Param("category")
	var categoryTasks []Task
	for _, task := range tasks {
		if strings.EqualFold(task.Category, category) {
			categoryTasks = append(categoryTasks, task)
		}
	}
	c.JSON(http.StatusOK, categoryTasks)
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

	for i, task := range tasks {
		if task.ID == id {
			for _, user := range tasks[i].SharedWith {
				if user == req.User {
					c.JSON(http.StatusConflict, gin.H{"error": "Task already shared with this user"})
					return
				}
			}
			tasks[i].SharedWith = append(tasks[i].SharedWith, req.User)
			c.JSON(http.StatusOK, tasks[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func getSharedTasks(c *gin.Context) {
	user := c.Param("user")
	var shared []Task
	for _, task := range tasks {
		for _, u := range task.SharedWith {
			if strings.EqualFold(u, user) {
				shared = append(shared, task)
				break
			}
		}
	}
	if len(shared) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks shared with this user"})
	}
	c.JSON(http.StatusOK, shared)
}

func getRecentTasks(c *gin.Context) {
	minutesStr := c.Query("minutes")
	if minutesStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minutes query parameter is required"})
		return
	}

	minutes, err := strconv.Atoi(minutesStr)
	if err != nil || minutes <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minutes parameter"})
		return
	}

	timeLimit := time.Now().Add(-time.Duration(minutes) * time.Minute)
	var recentTasks []Task
	for _, task := range tasks {
		if task.CreatedAt.After(timeLimit) {
			recentTasks = append(recentTasks, task)
		}
	}
	if len(recentTasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No recent tasks found"})
		return
	}
	c.JSON(http.StatusOK, recentTasks)
}

func getTasksByOwner(c *gin.Context) {
	owner := c.Param("owner")
	var ownerTasks []Task
	for _, task := range tasks {
		if strings.EqualFold(task.Owner, owner) {
			ownerTasks = append(ownerTasks, task)
		}
	}
	if len(ownerTasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks found for this owner"})
		return
	}
	c.JSON(http.StatusOK, ownerTasks)
}

func getNextID() int64 {
	var maxID int64 = 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}
