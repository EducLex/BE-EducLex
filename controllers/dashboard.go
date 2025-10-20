package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Dashboard data count
func GetDashboardStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	articleCount, _ := config.ArticleCollection.CountDocuments(ctx, bson.M{})
	questionCount, _ := config.QuestionCollection.CountDocuments(ctx, bson.M{})
	tulisanCount, _ := config.TulisanCollection.CountDocuments(ctx, bson.M{})
	peraturanCount, _ := config.PeraturanCollection.CountDocuments(ctx, bson.M{})
	userCount, _ := config.UserCollection.CountDocuments(ctx, bson.M{})

	c.JSON(http.StatusOK, gin.H{
		"articles":   articleCount,
		"questions":  questionCount,
		"tulisan":    tulisanCount,
		"peraturan":  peraturanCount,
		"users":      userCount,
	})
}

// Daftar semua pengguna
func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}
	defer cursor.Close(ctx)

	var users []bson.M
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data pengguna"})
		return
	}

	c.JSON(http.StatusOK, users)
}
