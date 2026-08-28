package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"time"

	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RefreshToken(userService *service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("refresh_token")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "refresh token is required",
			})
			return
		}

		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])

		checkingToken, err := userService.VerfiyToken(ctx.Request.Context(), tokenHash)
		if err != nil {
			if errors.Is(err, service.ErrTokenNotFound) {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid or expired refresh token",
				})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to verify refresh token",
			})
			return
		}

		secret, ok := os.LookupEnv("JWT_SECRET")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error: JWT secret not configured",
			})
			return
		}

		claims := jwt.MapClaims{
			"user_id": checkingToken.ID,
			"email":   checkingToken.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
			"iat":     time.Now().Unix(),
		}

		accessToken := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			claims,
		)

		signedToken, err := accessToken.SignedString([]byte(secret))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate access token",
			})
			return
		}

		ctx.SetCookie(
			"access_token",
			signedToken,
			60*60*24,
			"/",
			"",
			false,
			true,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "token refreshed successfully",
			"token":   signedToken,
			"user":    checkingToken,
		})
	}
}
