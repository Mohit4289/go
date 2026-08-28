package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

func Logout(UserService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(500, gin.H{
				"message": "not able to fetch  refresh token",
				"err":     err,
			})
			return
		}

		hash := sha256.Sum256([]byte(refreshToken))
		tokenHash := hex.EncodeToString(hash[:])

		var validationErr service.ValidationError

		removingToken, err := UserService.LogoutRemoveToken(c, tokenHash)
		if err != nil {
			if errors.As(err, &validationErr) {
				c.JSON(500, gin.H{
					"field": validationErr.Field,
					"error": validationErr.Msg,
				})
				return
			}

			c.JSON(500, gin.H{
				"message": "failed to logout",
			})
			return
		}

		if !removingToken {
			c.JSON(404, gin.H{
				"message": "refresh token not found",
			})
			return
		}

		c.SetCookie(
			"access_token",
			"",
			-1,
			"/",
			"",
			true,
			true,
		)
		c.SetCookie(
			"refresh_token",
			"",
			-1,
			"/",
			"",
			true,
			true,
		)

		c.JSON(200, gin.H{
			"message": "successfully remove token and logout",
		})
	}
}
