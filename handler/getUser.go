package handler

import (
	"net/http"

	"gin-quickstart/repository"

	"github.com/gin-gonic/gin"
)

type GETUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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
		users := make([]GETUser, 0)

		rows, err := r.userRepo.DB.Query(
			ctx.Request.Context(),
			`SELECT id, name, email FROM public."user"`,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to query users",
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
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to scan user data",
				})
				return
			}

			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while reading users",
			})
			return
		}

		userID, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: user id not found in context",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"users":   users,
			"user_id": userID,
		})
	}
}
