package cache

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

// Connects to Redis and initializes the client
func Connect() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: "",
		DB:       0,
	})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Error to connect on Redis", err)
	}

	log.Println("Connected to Redis successfully")
}
