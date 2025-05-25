package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	// Get env variables with defaults
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	log.Printf("Connecting to Redis at %s", addr)

	dbClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := dbClient.Ping(context.TODO()).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis at %s: %v", addr, err)
		panic(err)
	}

	log.Printf("Successfully connected to Redis at %s", addr)
	Client = dbClient
}
