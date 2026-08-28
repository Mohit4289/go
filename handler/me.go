package handler

import (
	"errors"
	"net/http"

	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

func MeUser(u *service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idValue, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: user id not found in context",
			})
			return
		}

		id, ok := idValue.(int)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error: invalid user id format",
			})
			return
		}

		userData, err := u.FetchUser(ctx.Request.Context(), id)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "user not found",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch user data",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": userData,
		})
	}
}
