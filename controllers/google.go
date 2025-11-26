package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"fmt"
	"net/url"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

type GoogleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// --- STEP 1: Redirect ke Google ---
func GoogleLogin(c *gin.Context) {
    // ambil redirect dari FE kalau ada
    redirectURI := c.Query("redirect_uri")
    if redirectURI == "" {
        redirectURI = "http://localhost:8080/auth/google/callback" 
    }

    // override redirect URL di config
    config.GoogleOauthConfig.RedirectURL = redirectURI

    url := config.GoogleOauthConfig.AuthCodeURL(
        "state-token",
        oauth2.AccessTypeOffline,
        oauth2.ApprovalForce,
    )

    c.Redirect(http.StatusTemporaryRedirect, url)
}

// --- STEP: Callback (auto register or login) ---

func GoogleCallback(c *gin.Context) {
	// ambil "code" dari URL yang dikirim Google
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code tidak ada"})
		return
	}

	// tukar code -> token Google
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal tukar code ke token"})
		return
	}

	// pakai token untuk ambil data user Google
	client := config.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal ambil data user google"})
		return
	}
	defer resp.Body.Close()

	var gUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Id    string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal decode user google"})
		return
	}

	// cek di DB: kalau belum ada, buat user baru
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err = config.UserCollection.FindOne(ctx, bson.M{"email": gUser.Email}).Decode(&user)
	if err != nil {
		// kalau tidak ketemu → insert user baru
		user = models.User{
			ID:       primitive.NewObjectID(),
			Username: gUser.Name,
			Email:    gUser.Email,
			Role:     "user",
		}
		if _, err := config.UserCollection.InsertOne(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal simpan user google"})
			return
		}
	}
	// --- Buat JWT
	jwtToken, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username, user.Role)

	// Redirect ke frontend (index.html) dengan query params
	redirectURL := fmt.Sprintf(
		"http://127.0.0.1:5500/index.html?token=%s&user_id=%s&username=%s&email=%s&role=%s",
		url.QueryEscape(jwtToken),
		url.QueryEscape(user.ID.Hex()),
		url.QueryEscape(user.Username),
		url.QueryEscape(user.Email),
		url.QueryEscape(user.Role),
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)

}
