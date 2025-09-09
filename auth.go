package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// START login Google
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("random-state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback dari Google
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	tok, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var g map[string]any
	_ = json.Unmarshal(body, &g)

	email, _ := g["email"].(string)
	name, _ := g["name"].(string)
	if email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No email from Google"})
		return
	}

	// Simpan user (insert tanpa cek duplikat)
	_, _ = userCollection.InsertOne(context.Background(), bson.M{
		"email": email,
		"name":  name,
	})

	// Buat JWT
	jwtToken, err := GenerateJWT(email, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"email":   email,
		"name":    name,
		"token":   jwtToken,
	})
}

// END login Google

// Contoh route proteksi
func ProfileHandler(c *gin.Context) {
	email := c.GetString("email")
	name := c.GetString("name")
	c.JSON(http.StatusOK, gin.H{
		"email": email,
		"name":  name,
	})
}
