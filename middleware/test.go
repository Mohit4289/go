package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("Middleware: before")

		// Check
		name := c.GetHeader("name")

		if name == "" {
			c.JSON(400, gin.H{
				"message": "name mat bhejo blank",
			})
			c.Abort()
			return
		}

		if name != "mohit" {
			c.JSON(400, gin.H{
				"message": "wrong name",
			})
			c.Abort()
			return
		}

		// Name correct → continue to handler
		c.Next()

		fmt.Println("Middleware: after")
	}
}