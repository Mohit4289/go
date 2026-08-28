package handler

import (
	"errors"
	"net/http"

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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body: please provide email and pass",
			})
			return
		}

		token, refreshToken, err := userService.CheckPassword(
			c.Request.Context(),
			login.Email,
			login.Password,
		)
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid email or password",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		c.SetCookie(
			"refresh_token",
			refreshToken,
			60*60*24,
			"/",
			"",
			false,
			true,
		)

		c.SetCookie(
			"access_token",
			token,
			60*60*24,
			"/",
			"",
			false,
			true,
		)

		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
			"token":   token,
		})
	}
}
