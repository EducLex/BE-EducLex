package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Koleksi MongoDB
var articleCollection = config.ArticleCollection

// ✅ Tambah artikel baru (Admin)
func CreateArticle(c *gin.Context) {
	os.MkdirAll("uploads", os.ModePerm)

	var input models.Article
	input.Judul = c.PostForm("judul")
	input.Isi = c.PostForm("isi")
	input.Penulis = "Admin"
	input.CreatedAt = time.Now()

	file, _ := c.FormFile("gambar")
	if file != nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err == nil {
			input.Gambar = path
		}
	}

	dokumen, _ := c.FormFile("dokumen")
	if dokumen != nil {
		path := "uploads/" + dokumen.Filename
		if err := c.SaveUploadedFile(dokumen, path); err == nil {
			input.Dokumen = path
		}
	}

	collection := config.ArticleCollection
	res, err := collection.InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Gagal menambahkan artikel",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Artikel berhasil ditambahkan",
		"id":      res.InsertedID,
	})
}

// ✅ Ambil semua artikel (User & Admin)
func GetArticles(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ✅ gunakan config.ArticleCollection, bukan articleCollection
	cursor, err := config.ArticleCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer cursor.Close(ctx)

	var articles []models.Article
	for cursor.Next(ctx) {
		var artikel models.Article
		if err := cursor.Decode(&artikel); err == nil {
			articles = append(articles, artikel)
		}
	}

	c.JSON(http.StatusOK, articles)
}

// ✅ Ambil artikel berdasarkan ID
func GetArticleByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var artikel models.Article
	err = articleCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&artikel)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, artikel)
}

// ✅ Update artikel (Admin)
func UpdateArticle(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input models.Article
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"judul":     input.Judul,
			"isi":       input.Isi,
			"updatedAt": time.Now(),
		},
	}

	// ✅ pakai config.ArticleCollection
	_, err = config.ArticleCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui artikel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artikel berhasil diperbarui"})
}

// ✅ Hapus artikel (Admin)
func DeleteArticle(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// ✅ Gunakan config.ArticleCollection
	res, err := config.ArticleCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus artikel"})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artikel berhasil dihapus"})
}
