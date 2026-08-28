package handler

import (
	"errors"
	"net/http"

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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body: please provide name, valid email, and pass",
			})
			return
		}

		err = userService.ValidateUser(c.Request.Context(), user.Email)
		if err != nil {
			if errors.Is(err, service.ErrUserAlreadyExists) {
				c.JSON(http.StatusConflict, gin.H{
					"error": "user with this email already exists",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to validate user",
			})
			return
		}

		hashPass, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to process password",
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create user",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "user created successfully",
			"user":    addingUser,
		})
	}
}
