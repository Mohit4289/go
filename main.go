package main

import (
	"context"
	"gin-quickstart/database"
	db "gin-quickstart/db/sqlc"
	"gin-quickstart/router"
	"gin-quickstart/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not loaded:", err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	connPool, err := database.Connect()
	if err != nil {
		log.Fatal("db failed: ", err)
	}

	redisClient := database.ConnectRedis()
	if redisClient != nil {
		defer redisClient.Close()
	}
	defer connPool.Close()

	queries := db.New(connPool)
	userService := service.NewUserService(queries, redisClient)
	propertyService := service.NewPropertyService(queries)

	router.SetupRoutes(r, userService)
	router.SetupPropertyRoutes(r, propertyService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("ERR: Server shutdown error: %v", err)
	}
}
