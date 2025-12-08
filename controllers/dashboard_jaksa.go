package controllers

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
)

// Dashboard Jaksa data count
func GetJaksaDashboardStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Statistik untuk Jaksa
	questionCount, _ := config.QuestionCollection.CountDocuments(ctx, bson.M{})
	tulisanCount, _ := config.TulisanCollection.CountDocuments(ctx, bson.M{})

	// Ambil statistik lainnya jika diperlukan
	// contoh: count pertanyaan yang belum dijawab
	unansweredQuestionCount, _ := config.QuestionCollection.CountDocuments(ctx, bson.M{"status": "Belum Dijawab"})

	// Ambil statistik tulisan Jaksa
	jaksaCount, _ := config.JaksaCollection.CountDocuments(ctx, bson.M{})

	c.JSON(http.StatusOK, gin.H{
		"questions":            questionCount,
		"tulisan":              tulisanCount,
		"unanswered_questions": unansweredQuestionCount,
		"jaksa_count":          jaksaCount,
	})
}

// Menampilkan Daftar Pertanyaan yang Belum Dijawab oleh Jaksa
func GetUnansweredQuestions(c *gin.Context) {
	questionCollection := config.QuestionCollection
	if questionCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi database belum siap"})
		return
	}

	// Ambil pertanyaan dengan status "Belum Dijawab"
	cursor, err := questionCollection.Find(context.Background(), bson.M{"status": "Belum Dijawab"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pertanyaan"})
		return
	}
	defer cursor.Close(context.Background())

	var questions []models.Question
	if err := cursor.All(context.Background(), &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data pertanyaan"})
		return
	}

	c.JSON(http.StatusOK, questions)
}
