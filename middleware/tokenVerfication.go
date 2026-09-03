package middleware

import (
	"net/http"

	"gin-quickstart/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenVerification() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required: access token is missing",
			})
			c.Abort()
			return
		}

		config := config.Load()
		secret := config.JWT_SECRET
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error: JWT secret not configured",
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		userIDVal, exists := claims["user_id"]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user_id not found in token claims",
			})
			c.Abort()
			return
		}

		var id int
		switch v := userIDVal.(type) {
		case float64:
			id = int(v)
		case int:
			id = v
		case int64:
			id = int(v)
		default:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id format in token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", id)
		c.Next()
	}
}
