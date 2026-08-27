package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"gin-quickstart/service"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RefreshToken(UserService *service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("refresh_token")
		if err != nil {
			ctx.JSON(401, gin.H{
				"message": "not able to get token",
				"err":     err,
			})
			return
		}

		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])

		checkingToken, err := UserService.VerfiyToken(ctx, tokenHash)
		if err != nil {
			ctx.JSON(401, gin.H{
				"message": "not verfied",
				"err":     err,
			})
			return
		}

		secret, ok := os.LookupEnv("JWT_SECRET")
		if !ok {
			ctx.JSON(401, gin.H{
				"message": "Token not configuered ",
			})
			return
		}

		claims := jwt.MapClaims{
			"user_id": checkingToken.ID,
			"email":   checkingToken.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
			"iat":     time.Now().Unix(),
		}

		access_token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			claims,
		)

		signedToken, err := access_token.SignedString([]byte(secret))
		if err != nil {
			ctx.JSON(401, gin.H{
				"message": "not able to generate access token ",
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

		ctx.JSON(200, gin.H{
			"message": "verfied succesfully",
			"data":    checkingToken,
		})

	}
}
