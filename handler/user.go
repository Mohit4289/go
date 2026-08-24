package handler

import (
	"gin-quickstart/service"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
)

type CreateUser struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"pass" binding:"required"`
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

		err = userService.ValidateUser(c.Request.Context(), user.Email)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		hashPass, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "hashpass error",
			})
			return
		}

		addingUser, err := userService.AddUser(
			c.Request.Context(),
			user.Name,
			user.Email,
			hashPass,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"message": "user created",
			"user":    addingUser,
		})
	}
}
