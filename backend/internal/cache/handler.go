package cache

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Registers cache-related routes in Gin
func RegisterRoutes(r *gin.Engine) {
	r.POST("/cache/set", setCache)
	r.GET("/cache/get/:key", getCache)
}

// Sets a key-value pair in Redis cache
func setCache(c *gin.Context) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := RDB.Set(Ctx, req.Key, req.Value, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set cache"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cache set successfully"})
}

// Retrieves a value from Redis cache by key
func getCache(c *gin.Context) {
	key := c.Param("key")
	value, err := RDB.Get(Ctx, key).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}
