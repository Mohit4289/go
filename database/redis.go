package database

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		DialTimeout: 1 * time.Second,
		MaxRetries:  1,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("WARN: Redis connection failed (running without Redis cache):", err)
		return nil
	}

	log.Println("INFO: Redis connected")
	return client
}
