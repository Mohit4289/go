package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := client.Ping(context.Background()).Err()

	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected")
	return client
}
