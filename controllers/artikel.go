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

func CreateArticle(c *gin.Context) {
	// Membuat folder uploads jika belum ada
	os.MkdirAll("uploads", os.ModePerm)

	var input models.Article
	// Mengambil data dari Form Data (bukan JSON)
	input.Judul = c.PostForm("judul")
	input.Isi = c.PostForm("isi")

	// Mengambil categoryId dari Form Data (sebagai string)
	categoryID := c.PostForm("categoryId")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID tidak boleh kosong"})
		return
	}

	// Mengonversi categoryId menjadi primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID tidak valid"})
		return
	}
	input.CategoryID = objectID

	// Menangani file gambar
	file, _ := c.FormFile("gambar")
	if file != nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err == nil {
			input.Gambar = path
		}
	}

	// Menangani file dokumen
	dokumen, _ := c.FormFile("dokumen")
	if dokumen != nil {
		path := "uploads/" + dokumen.Filename
		if err := c.SaveUploadedFile(dokumen, path); err == nil {
			input.Dokumen = path
		}
	}

	// Insert artikel ke MongoDB
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

// ✅ Ambil semua artikel (User & Admin) berdasarkan CategoryID
func GetArticles(c *gin.Context) {
	categoryID := c.DefaultQuery("categoryId", "")

	// Jika categoryId ada, gunakan filter
	var filter bson.M
	if categoryID != "" {
		// Validasi categoryId
		_, err := primitive.ObjectIDFromHex(categoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID tidak valid"})
			return
		}

		filter = bson.M{"categoryId": categoryID}
	} else {
		filter = bson.M{} 
	}

	cursor, err := config.ArticleCollection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var articles []models.Article
	if err := cursor.All(context.Background(), &articles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
	// Ambil ID artikel
	idParam := c.Param("id")
	articleID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID artikel tidak valid"})
		return
	}

	// Ambil data form
	judul := c.PostForm("judul")
	isi := c.PostForm("isi")
	categoryIDStr := c.PostForm("categoryId")
	penulis := c.PostForm("penulis")

	if judul == "" || isi == "" || categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Judul, isi, dan category wajib diisi",
		})
		return
	}

	// Convert categoryId ke ObjectID
	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID tidak valid"})
		return
	}

	// Ambil category dari tabel category
	var category models.Category
	err = config.CategoryCollection.FindOne(
		context.Background(),
		bson.M{"_id": categoryID},
	).Decode(&category)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category tidak ditemukan"})
		return
	}

	// Siapkan data update
	updateData := bson.M{
		"judul":         judul,
		"isi":           isi,
		"category_id":   categoryID,
		"category_nama": category.Name,
		"penulis":       penulis,
		"updatedAt":     time.Now(),
	}

	// Handle dokumen (optional)
	dokumen, err := c.FormFile("dokumen")
	if err == nil {
		path := "uploads/" + dokumen.Filename
		if err := c.SaveUploadedFile(dokumen, path); err == nil {
			updateData["dokumen"] = path
		}
	}

	// Update ke MongoDB
	result, err := config.ArticleCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": articleID},
		bson.M{"$set": updateData},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update artikel"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artikel tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Artikel berhasil diperbarui",
	})
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
