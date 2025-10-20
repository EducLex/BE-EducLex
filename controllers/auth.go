package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, _ := config.UserCollection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"email": input.Email},
			{"username": input.Username},
		},
	})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Email already exists"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}
	_, err := config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{"message": "register success", "token": token})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := config.UserCollection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{"message": "login success", "token": token})
}

// Register Admin
func RegisterAdmin(c *gin.Context) {
	var input struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, _ := config.UserCollection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"email": input.Email},
			{"username": input.Username},
		},
	})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Email already exists"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "admin", // <- langsung admin
	}
	_, err := config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}

	token, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{"message": "admin register success", "token": token})
}

// Logout
func Logout(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")
	if tokenHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header tidak ditemukan"})
		return
	}

	token := tokenHeader
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blacklist := models.TokenBlacklist{
		Token:     token,
		ExpiredAt: time.Now().Add(24 * time.Hour), // token dianggap kadaluarsa 1 hari
	}

	_, err := config.TokenBlacklistCollection.InsertOne(ctx, blacklist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan token ke blacklist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}
