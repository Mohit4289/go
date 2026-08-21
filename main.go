package main

import (
	"gin-quickstart/database"
	"gin-quickstart/repository"
	"gin-quickstart/router"
	"gin-quickstart/service"
	"log"

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

	router.SetupRoutes(r, userService)

	r.Run()
}
