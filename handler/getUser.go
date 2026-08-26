package handler

import (
	"gin-quickstart/repository"

	"github.com/gin-gonic/gin"
)

type GETUser struct {
	ID    int64
	Name  string
	Email string
}

type UserService struct {
	userRepo *repository.UserRepo
}

func GetUserRepo(userRepo *repository.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func GetUser(r *UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var users []GETUser

		rows, err := r.userRepo.DB.Query(
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

			var user GETUser

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

		user_data, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(401, gin.H{
				"message": "data didnt came in main",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"users":     users,
			"user_data": user_data,
		})
	}
}
