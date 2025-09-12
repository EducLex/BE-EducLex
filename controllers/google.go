package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoogleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func GoogleLoginRedirect(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL(uuid.NewString())
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

	client := config.GoogleOauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info"})
		return
	}

	var gUser GoogleUser
	if err := json.Unmarshal(data, &gUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err = config.UserCollection.FindOne(ctx, bson.M{"google_id": gUser.ID}).Decode(&user)

	// Jika user belum ada → buat user baru
	if err != nil {
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
	}

	// Buat JWT internal untuk aplikasi
	jwtToken, err := middleware.GenerateJWT(user.ID.Hex(), user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Simpan token ke DB
	_, err = config.UserCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"token": jwtToken}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	// Redirect ke frontend dengan membawa token
	redirectURL := os.Getenv("FRONTEND_URL") + "/google-success?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
