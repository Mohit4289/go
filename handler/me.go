package handler

import (
	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

func MeUser(u *service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idValue, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(401, gin.H{
				"message": "user didn't came in",
			})
			return
		}

		id, ok := idValue.(int8)
		if !ok {
			ctx.JSON(500, gin.H{
				"message": "invalid user_id type",
			})
			return
		}

		userData, err := u.FetchUser(ctx, id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "not fetched maybe user dont exsit",
				"error":   err,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"data": userData.Name,
		})
		return
	}
}
