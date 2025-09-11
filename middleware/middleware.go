package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Format header: "Bearer <token>"
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Parse JWT
		secret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// simpan claims ke context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user", claims)
		}

		c.Next()
	}
}
