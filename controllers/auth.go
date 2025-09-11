package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("SECRET_KEY_KAMU")

// --- REGISTER MANUAL ---
func Register(c *gin.Context) {
	var req struct {
		Email           string `json:"email"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// cek apakah email sudah ada
	var existing models.User
	err := config.UserCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	newUser := models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashed),
		Provider:  "manual",
		CreatedAt: time.Now(),
	}
	_, err = config.UserCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// --- LOGIN MANUAL ---
func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	err := config.UserCollection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Kalau user dari Google → tolak login manual
	if user.Provider == "google" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This account uses Google login"})
		return
	}

	// Cek password manual
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Buat JWT
	claims := jwt.MapClaims{
		"email": user.Email,
		"name":  user.Username,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{
		"token": jwtString,
		"user":  user.Email,
	})
}
