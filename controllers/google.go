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
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
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

	data, _ := ioutil.ReadAll(resp.Body)
	var gUser GoogleUser
	json.Unmarshal(data, &gUser)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err = config.UserCollection.FindOne(ctx, bson.M{"google_id": gUser.ID}).Decode(&user)
	if err != nil {
		// buat user baru
		user = models.User{
			ID:       primitive.NewObjectID(),
			Username: gUser.Name + "-" + uuid.NewString()[:6],
			Email:    gUser.Email,
			GoogleID: gUser.ID,
		}
		config.UserCollection.InsertOne(ctx, user)
	}

	jwtToken, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username)
	redirectURL := os.Getenv("FRONTEND_URL") + "/google-success?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
