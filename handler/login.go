package handler

import (
	"errors"
	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"pass" binding:"required"`
}

func Login(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var login LoginUser

		err := c.ShouldBindJSON(&login)

		if err != nil {
			c.JSON(400, gin.H{
				"message": "fill the data",
			})
			return
		}

		var validationErr service.ValidationError

		token, err := userService.CheckPassword(
			c.Request.Context(),
			login.Email,
			login.Password,
		)
		if err != nil {
			if errors.As(err, &validationErr) {
				c.JSON(400, gin.H{
					"field": validationErr.Field,
					"error": validationErr.Msg,
				})
				return
			}

			c.JSON(500, gin.H{
				"error": "internal server error",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "login sucessfully",
			"token":   token,
		})
	}
}
