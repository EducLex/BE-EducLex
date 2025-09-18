package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET semua artikel
func GetArticles(c *gin.Context) {
	cursor, err := config.ArticleCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer cursor.Close(context.Background())

	var articles []models.Article
	if err := cursor.All(context.Background(), &articles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal decode data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": articles})
}

// GET artikel by ID
func GetArticleByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var article models.Article
	err = config.ArticleCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&article)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// POST tambah artikel
func CreateArticle(c *gin.Context) {
	var newArticle models.Article

	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newArticle.ID = primitive.NewObjectID()
	newArticle.CreatedAt = time.Now()

	_, err := config.ArticleCollection.InsertOne(context.Background(), newArticle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan artikel"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Artikel berhasil ditambahkan", "data": newArticle})
}

// PUT update artikel
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var updateData models.Article
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":      updateData.Title,
			"content":    updateData.Content,
			"image":      updateData.Image,
			"file":       updateData.File,
			"created_at": time.Now(),
		},
	}

	_, err = config.ArticleCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update artikel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artikel berhasil diupdate"})
}

// DELETE artikel
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	_, err = config.ArticleCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus artikel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artikel berhasil dihapus"})
}
