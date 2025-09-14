package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

type GoogleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// --- STEP 1: Redirect ke Google ---
func GoogleLogin(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// --- STEP 2: Register pakai Google ---
func GoogleRegister(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := config.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var gUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// cek apakah email/GoogleID sudah ada
	count, _ := config.UserCollection.CountDocuments(ctx, bson.M{"google_id": gUser.ID})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Google account already registered"})
		return
	}

	// simpan user baru
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: gUser.Name,
		Email:    gUser.Email,
		GoogleID: gUser.ID,
	}
	_, err = config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Google register success"})
}

// --- STEP 3: Login pakai Google (hanya jika sudah register) ---
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := config.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var gUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// cek apakah user sudah ada
	var user models.User
	err = config.UserCollection.FindOne(ctx, bson.M{"google_id": gUser.ID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Google account not registered. Please register first."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// generate JWT
	jwtToken, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "Google login success",
		"token":   jwtToken,
	})

}
