package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProfileHandler -> hanya bisa diakses kalau JWT valid
func ProfileHandler(c *gin.Context) {
	// Ambil data user dari context yang di-set middleware
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile data",
		"user":    userData,
	})
}
