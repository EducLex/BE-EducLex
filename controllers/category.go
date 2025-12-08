package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Fungsi untuk seed kategori default
func SeedCategories() {
	categories := []models.Category{
		{Name: "Pembinaan"},
		{Name: "Intelijen"},
		{Name: "Pidana Umum"},
		{Name: "Pidana Khusus"},
		{Name: "Perdata dan Tata Usaha Negara"},
		{Name: "Pidana Militer"},
		{Name: "Asisten Pengawasan"},
	}

	for _, category := range categories {
		// Cek jika kategori sudah ada
		var existing models.Category
		err := config.CategoryCollection.FindOne(context.Background(), bson.M{"name": category.Name}).Decode(&existing)
		if err != nil && err.Error() != "mongo: no documents in result" {
			log.Fatalf("Failed to check category: %v", err)
		}

		// Jika kategori belum ada, masukkan ke dalam koleksi
		if existing.Name == "" {
			category.ID = primitive.NewObjectID()
			_, err = config.CategoryCollection.InsertOne(context.Background(), category)
			if err != nil {
				log.Fatalf("Failed to insert category: %v", err)
			}
		}
	}
}

// Fungsi untuk menambahkan kategori baru
func CreateCategory(c *gin.Context) {
	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menambahkan kategori", "detail": err.Error()})
		return
	}

	// Validasi kategori dan subkategori
	if input.Name != "internal" && input.Name != "eksternal" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori harus 'internal' atau 'eksternal'"})
		return
	}

	// Validasi subkategori internal
	if input.Name == "internal" {
		internalCategories := []string{"Pembinaan", "Intelijen", "Pidana Umum", "Pidana Khusus", "Perdata dan Tata Usaha Negara", "Pidana Militer", "Asisten Pengawasan"}
		subCategory := c.DefaultPostForm("subkategori", "")
		validSubCategory := false
		for _, category := range internalCategories {
			if subCategory == category {
				validSubCategory = true
				break
			}
		}
		if !validSubCategory {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subkategori internal tidak valid"})
			return
		}
	}

	// Validasi subkategori eksternal
	if input.Name == "eksternal" {
		eksternalCategories := []string{"Peraturan UUD", "Peraturan Presiden", "Perpres", "Keppres"}
		subCategory := c.DefaultPostForm("subkategori", "")
		validSubCategory := false
		for _, category := range eksternalCategories {
			if subCategory == category {
				validSubCategory = true
				break
			}
		}
		if !validSubCategory {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subkategori eksternal tidak valid"})
			return
		}
	}

	// Menyimpan data kategori baru
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	input.ID = primitive.NewObjectID()

	_, err := config.CategoryCollection.InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan kategori", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil ditambahkan",
		"data":    input,
	})
}

// Fungsi untuk mengambil kategori berdasarkan ID
func GetCategoryByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var category models.Category
	err = config.CategoryCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&category)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Fungsi untuk mengambil semua kategori
func GetCategories(c *gin.Context) {
	cursor, err := config.CategoryCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var categories []models.Category
	if err := cursor.All(context.Background(), &categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Update kategori
func UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid", "detail": err.Error()})
		return
	}

	// Periksa apakah collection sudah diinisialisasi
	if config.CategoryCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Category collection not initialized"})
		return
	}

	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca input", "detail": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":      input.Name,
			"updatedAt": time.Now(),
		},
	}

	result, err := config.CategoryCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kategori", "detail": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil diperbarui"})
}

// Hapus kategori
func DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid", "detail": err.Error()})
		return
	}

	// Periksa apakah collection sudah diinisialisasi
	if config.CategoryCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Category collection not initialized"})
		return
	}

	res, err := config.CategoryCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori", "detail": err.Error()})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}
