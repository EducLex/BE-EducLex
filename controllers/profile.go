package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProfileHandler -> hanya bisa diakses kalau JWT valid
func ProfileHandler(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"message":  "Profile data",
		"user_id":  userID,
		"username": username,
		"role":     role,
	})
}
