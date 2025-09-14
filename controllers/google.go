package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
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

// --- STEP 1: Redirect ke Google untuk Login
func GoogleLogin(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL("login", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// --- STEP 1: Redirect ke Google untuk Register
func GoogleRegister(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL("register", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// --- STEP 2: Callback setelah pilih akun Google
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// Tukar code → access token
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Ambil data user dari Google
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

	var user models.User
	err = config.UserCollection.FindOne(ctx, bson.M{"google_id": gUser.ID}).Decode(&user)

	if state == "register" {
		// Mode Register
		if err == mongo.ErrNoDocuments {
			user = models.User{
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
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already registered"})
			return
		}
	} else if state == "login" {
		// Mode Login
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not registered"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	// --- STEP 3: Buat JWT
	jwtToken, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username)

	// --- STEP 4: Redirect ke Frontend dengan token
	redirectURL := os.Getenv("FRONTEND_URL") + "/google-success?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
