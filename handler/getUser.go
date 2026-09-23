package handler

import (
	"net/http"

	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

func GetUser(userService *service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := userService.ListUsers(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to query users",
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

