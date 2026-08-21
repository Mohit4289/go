package handler

import (
	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

type CreateUser struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,min=18"`
}

func User(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user CreateUser

		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "invalid request",
			})
			return
		}

		err = userService.ValidateUser(user.Email)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "working",
			"user":    user,
		})
	}
}
