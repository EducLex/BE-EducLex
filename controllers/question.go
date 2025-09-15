package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tambah pertanyaan
func CreateQuestion(c *gin.Context) {
	var input struct {
		Nama       string `json:"nama"`
		Email      string `json:"email"`
		Pertanyaan string `json:"pertanyaan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question := models.Question{
		ID:         primitive.NewObjectID(),
		Nama:       input.Nama,
		Email:      input.Email,
		Pertanyaan: input.Pertanyaan,
		Jawaban:    "Pertanyaanmu sedang diproses oleh Jaksa EducLex...",
		CreatedAt:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.QuestionCollection.InsertOne(ctx, question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pertanyaan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pertanyaan berhasil dikirim", "data": question})
}

// Ambil semua pertanyaan
func GetQuestions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.QuestionCollection.Find(ctx, primitive.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer cursor.Close(ctx)

	var questions []models.Question
	if err := cursor.All(ctx, &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal decode data"})
		return
	}

	c.JSON(http.StatusOK, questions)
}
