package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenVerification() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "authentication required",
			})
			c.Abort()
			return
		}

		secret, ok := os.LookupEnv("JWT_SECRET")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Secret token not configured",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{
				"message": "user id not fetched",
			})
			c.Abort()
			return
		}

		user_id := claims["user_id"]
		id := int(user_id.(float64))
		c.Set("user_id", id)
		c.Next()

	}
}
