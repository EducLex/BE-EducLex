package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Koneksi MongoDB
	ConnectDB()

	// Init router
	r := gin.Default()

	// Routes Google OAuth
	r.GET("/auth/google/login", GoogleLogin)       // redirect ke Google
	r.GET("/auth/google/callback", GoogleCallback) // callback dari Google

	// Contoh route proteksi JWT
	r.GET("/profile", AuthMiddleware(), ProfileHandler)

	// Root test
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Educlex API 🚀"})
	})

	// Jalankan server
	r.Run(":8080")
}
