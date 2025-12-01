package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

var jwtKey = []byte("superSecretKey123") 

// GenerateJWT buat token JWT (dengan role)
func GenerateJWT(userID string, username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
	
}

// AuthMiddleware verifikasi JWT + cek blacklist
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 🔒 Cek apakah token sudah di-blacklist
		if IsTokenBlacklisted(tokenString) {
			c.JSON(401, gin.H{"error": "Token sudah tidak berlaku (logout)"})
			c.Abort()
			return
		}

		// Parse JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		// simpan ke context
		c.Set("user_id", fmt.Sprintf("%v", claims["user_id"]))
		c.Set("username", fmt.Sprintf("%v", claims["username"]))
		c.Set("role", fmt.Sprintf("%v", claims["role"]))
		c.Next()
	}
}

// AdminMiddleware -> hanya admin yang bisa akses
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(403, gin.H{"error": "Forbidden: Admins only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 🔍 Fungsi bantu: cek apakah token ada di blacklist
func IsTokenBlacklisted(token string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, _ := config.TokenBlacklistCollection.CountDocuments(ctx, bson.M{"token": token})
	return count > 0
}
