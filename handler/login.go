package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=16"`
}

func Login(c *gin.Context) {
	fmt.Println("login api")

	var login LoginUser

	err := c.ShouldBindJSON(&login)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "fill the data",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "login sucessfully",
		"user":    login,
	})

}
