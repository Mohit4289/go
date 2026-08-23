package handler

import (
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

		checkingpass, err := userService.CheckPassword(login.Email, login.Password)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "login sucessfully",
			"user":    login,
			"bool":    checkingpass,
		})
	}
}
