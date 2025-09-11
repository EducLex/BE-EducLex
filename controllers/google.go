package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/idtoken"
)

var (
	googleClientID     = "YOUR_CLIENT_ID"
	googleClientSecret = "YOUR_CLIENT_SECRET"
	redirectURI        = "http://localhost:8080/auth/google/callback"
)

// --- Redirect ke Google ---
func GoogleLogin(c *gin.Context) {
	authURL := "https://accounts.google.com/o/oauth2/v2/auth"
	params := url.Values{}
	params.Add("client_id", googleClientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "openid email profile")
	params.Add("access_type", "offline")
	params.Add("prompt", "select_account")

	c.Redirect(http.StatusTemporaryRedirect, authURL+"?"+params.Encode())
}

// --- Callback dari Google ---
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// tukar code dengan token
	tokenURL := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", googleClientID)
	data.Set("client_secret", googleClientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token"})
		return
	}
	defer resp.Body.Close()

	var tokenResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token"})
		return
	}

	idToken, ok := tokenResp["id_token"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token received"})
		return
	}

	// verifikasi id_token
	payload, err := idtoken.Validate(c, idToken, googleClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID Token"})
		return
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)

	// cek apakah user sudah ada
	var user models.User
	err = config.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)

	if err != nil { // kalau belum ada → daftar baru
		newUser := models.User{
			Email:     email,
			Username:  name,
			Password:  "", // kosong karena pakai Google
			Provider:  "google",
			CreatedAt: time.Now(),
		}
		_, err := config.UserCollection.InsertOne(context.Background(), newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
			return
		}
		user = newUser
	}

	// generate JWT
	claims := jwt.MapClaims{
		"email": user.Email,
		"name":  user.Username,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{
		"token": jwtString,
		"user":  user,
	})
}
