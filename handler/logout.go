package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

func Logout(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "refresh token is required to logout",
			})
			return
		}

		hash := sha256.Sum256([]byte(refreshToken))
		tokenHash := hex.EncodeToString(hash[:])

		removingToken, err := userService.LogoutRemoveToken(c.Request.Context(), tokenHash)
		if err != nil {
			if errors.Is(err, service.ErrTokenNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "refresh token not found or already invalidated",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to logout",
			})
			return
		}

		if !removingToken {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "refresh token not found",
			})
			return
		}

		c.SetCookie(
			"access_token",
			"",
			-1,
			"/",
			"",
			false,
			true,
		)
		c.SetCookie(
			"refresh_token",
			"",
			-1,
			"/",
			"",
			false,
			true,
		)

		c.JSON(http.StatusOK, gin.H{
			"message": "logged out successfully",
		})
	}
}
