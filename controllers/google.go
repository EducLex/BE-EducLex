package controllers

import (
	"context"
	"encoding/json"
	"io" // ganti ioutil
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

func GoogleLogin(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline, // biar dapat refresh token
		oauth2.ApprovalForce,     // paksa Google kasih refresh token
	)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// Tukar code dari Google dengan access token
	token, err := config.GoogleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Pakai access token untuk ambil data user Google
	client := config.GoogleOauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var gUser GoogleUser
	if err := json.Unmarshal(data, &gUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Cari user berdasarkan google_id
	var user models.User
	err = config.UserCollection.FindOne(ctx, bson.M{"google_id": gUser.ID}).Decode(&user)

	// Jika user belum ada, buat baru
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
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Buat JWT untuk aplikasi
	jwtToken, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username)

	// Update token di database
	_, err = config.UserCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"token": jwtToken}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	// Redirect ke frontend dengan JWT
	redirectURL := os.Getenv("FRONTEND_URL") + "/google-success?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
