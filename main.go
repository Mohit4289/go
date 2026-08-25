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

type User struct {
	ID    int64
	Name  string
	Email string
}

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

	r.GET("/user", func(ctx *gin.Context) {

		var users []User

		rows, err := db.Query(
			ctx,
			`SELECT id, name, email FROM public."user"`,
		)

		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "db issue",
			})
			return
		}

		defer rows.Close()

		for rows.Next() {

			var user User

			err := rows.Scan(
				&user.ID,
				&user.Name,
				&user.Email,
			)

			if err != nil {
				ctx.JSON(500, gin.H{
					"message": "failed to read user",
				})
				return
			}

			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			ctx.JSON(500, gin.H{
				"message": "failed while reading users",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"users": users,
		})
	})

	r.Run()
}
