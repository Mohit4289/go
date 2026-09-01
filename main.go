package main

import (
	"context"
	"gin-quickstart/database"
	"gin-quickstart/handler"
	"gin-quickstart/repository"
	"gin-quickstart/router"
	"gin-quickstart/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env:", err)
	}

	r := gin.Default()

	db, err := database.Connect()
	if err != nil {
		log.Fatal("db failed", err)
	}

	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userData := handler.GetUserRepo(userRepo)

	router.SetupRoutes(r, userService, userData)

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
		log.Printf("Server shutdown error: %v", err)
	}
}
