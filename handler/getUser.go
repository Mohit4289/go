package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GETUser struct {
	ID    int64
	Name  string
	Email string
}

type UserRepo struct {
	db *pgxpool.Pool
}

func GetUserRepo(DB *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: DB,
	}
}

func GetUser(r *UserRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var users []GETUser

		rows, err := r.db.Query(
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
