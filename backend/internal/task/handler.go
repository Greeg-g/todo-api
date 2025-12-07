package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Greeg-g/todo-api/internal/auth"
	"github.com/Greeg-g/todo-api/internal/cache"
	"github.com/Greeg-g/todo-api/internal/database"
	"github.com/Greeg-g/todo-api/internal/model"
	"github.com/gin-gonic/gin"
)

var tasks = []model.Task{}

// RegisterRoutes registers task routes and applies JWT middleware to the group.
func RegisterRoutes(r *gin.Engine) {
	taskGroup := r.Group("/tasks")
	taskGroup.Use(auth.JWTMiddleware())
	{
		taskGroup.GET("/", getAllTasks)
		taskGroup.POST("/create", createTask)
		taskGroup.POST("/complete/:id", completeTask)
		taskGroup.DELETE("/delete/:id", deleteTask)
		taskGroup.GET("/category/:category", getCategoryTasks)
		taskGroup.POST("/share/:id", shareTask)
		taskGroup.GET("/shared", getSharedTasks)
		taskGroup.GET("/recent", getRecentTasks)
		taskGroup.GET("/owner/:owner", getTasksByOwner)
		taskGroup.GET("check-deadlines", checkDeadlines)
	}
}

// getAllTasks returns tasks owned by the authenticated user.
func getAllTasks(c *gin.Context) {
	u, ok := c.Get("currentUser")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := u.(*model.User)

	var tasks []model.Task
	if err := database.DB.Where("owner = ?", user.Username).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// createTask creates a new task owned by the authenticated user.
func createTask(c *gin.Context) {
	var newTask model.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask.Title = strings.TrimSpace(newTask.Title)
	if newTask.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}
	if len(newTask.Title) > 300 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title too long (max 300 characters)"})
		return
	}
	if len(newTask.Description) > 5000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description too long (max 5000 characters)"})
		return
	}

	if newTask.Deadline == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deadline is required"})
		return
	}

	u, ok := c.Get("currentUser")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := u.(*model.User)

	newTask.Owner = user.Username
	newTask.CreatedAt = time.Now()

	// Store original deadline for response (before adjustment).
	originalDeadline := newTask.Deadline

	// Adjust deadline +3 hours to compensate for timezone offset in database storage.
	adjustedDeadline := newTask.Deadline.Add(3 * time.Hour)
	if adjustedDeadline.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deadline must be a future date/time"})
		return
	}
	newTask.Deadline = &adjustedDeadline

	if err := database.DB.Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// Return original deadline (without +3h adjustment) to user.
	newTask.Deadline = originalDeadline
	c.JSON(http.StatusCreated, newTask)
}

// completeTask marks a task as completed; only the owner or a user in SharedWith
// is permitted to perform this action.
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

	u, ok := c.Get("currentUser")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := u.(*model.User)

	allowed := false
	if strings.EqualFold(task.Owner, user.Username) {
		allowed = true
	} else {
		for _, s := range task.SharedWith {
			if strings.EqualFold(s, user.Username) {
				allowed = true
				break
			}
		}
	}
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to complete this task"})
		return
	}

	task.Completed = true
	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// deleteTask removes a task by ID; only the owner is allowed to delete.
func deleteTask(c *gin.Context) {
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

	u, ok := c.Get("currentUser")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := u.(*model.User)

	if !strings.EqualFold(task.Owner, user.Username) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can delete this task"})
		return
	}

	if err := database.DB.Delete(&model.Task{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// getCategoryTasks returns tasks for the current user filtered by category.
func getCategoryTasks(c *gin.Context) {
	category := c.Param("category")

	if err := database.DB.Where("LOWER(category) = LOWER(?)", category).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// shareTask adds a target user (by username or email) to a task's SharedWith list.
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

	var target model.User
	userParam := strings.ToLower(strings.TrimSpace(req.User))
	if err := database.DB.Where("LOWER(username) = ? OR LOWER(email) = ?", userParam, userParam).First(&target).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target user not found"})
		return
	}

	if task.Owner == req.User {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Owner already has access to the task"})
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

// getSharedTasks returns tasks that have been shared with the authenticated user.
func getSharedTasks(c *gin.Context) {
	u, ok := c.Get("currentUser")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := u.(*model.User)
	var tasks []model.Task

	err := database.DB.Where("? = ANY (shared_with)", user.Username).Find(&tasks).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks shared with this user"})
	}

	c.JSON(http.StatusOK, tasks)
}

// getTasksByOwner returns tasks for the specified owner. If the requester is
// the owner all tasks are returned; otherwise only tasks shared with the
// requester are returned.
func getTasksByOwner(c *gin.Context) {
	owner := c.Param("owner")
	u, ok := c.Get("currentUser")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := u.(*model.User)

	var tasks []model.Task
	if err := database.DB.Where("LOWER(owner) = LOWER(?)", owner).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	if strings.EqualFold(owner, user.Username) {
		c.JSON(http.StatusOK, tasks)
		return
	}

	var filtered []model.Task
	for _, t := range tasks {
		for _, s := range t.SharedWith {
			if strings.EqualFold(s, user.Username) {
				filtered = append(filtered, t)
				break
			}
		}
	}

	if len(filtered) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to view these tasks"})
		return
	}
	c.JSON(http.StatusOK, filtered)
}

// getRecentTasks returns the most recent tasks for the authenticated user,
// using a per-user cache stored in Redis.
func getRecentTasks(c *gin.Context) {
	u, ok := c.Get("currentUser")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	user := u.(*model.User)

	key := fmt.Sprintf("recent_tasks:%d", user.ID)
	val, err := cache.RDB.Get(cache.Ctx, key).Result()
	if err == nil {
		c.Data(http.StatusOK, "application/json", []byte(val))
		return
	}

	var tasks []model.Task
	if err := database.DB.Where("owner = ?", user.Username).Order("created_at desc").Limit(10).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	jsonData, _ := json.Marshal(tasks)
	cache.RDB.Set(cache.Ctx, key, jsonData, time.Minute*10)

	c.JSON(http.StatusOK, tasks)
}

// checkDeadlines triggers an asynchronous enqueue of upcoming task deadlines.
func checkDeadlines(c *gin.Context) {
	go EnqueueUpcomingDeadlines()
	c.JSON(http.StatusOK, gin.H{"message": "Deadline check initiated"})
}
