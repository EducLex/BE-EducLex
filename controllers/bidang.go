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
	"go.mongodb.org/mongo-driver/mongo"
)

// Fungsi untuk menambahkan bidang baru
func CreateBidang(c *gin.Context) {
	var input models.Bidang
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menambahkan bidang", "detail": err.Error()})
		return
	}

	// Validasi status
	if input.Status != 0 && input.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus 0 (Non-Aktif) atau 1 (Aktif)"})
		return
	}

	// Set ID baru untuk bidang
	input.ID = primitive.NewObjectID()

	// Insert data bidang ke MongoDB
	_, err := config.BidangCollection.InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan bidang", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bidang berhasil ditambahkan",
		"id":      input.ID,
	})
}

// Ambil semua bidang
func GetBidangs(c *gin.Context) {
	// Periksa apakah collection sudah diinisialisasi
	if config.BidangCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bidang collection not initialized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.BidangCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bidang", "detail": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var bidangs []models.Bidang
	for cursor.Next(ctx) {
		var bidang models.Bidang
		if err := cursor.Decode(&bidang); err != nil {
			log.Printf("Error decoding bidang: %v", err)
			continue
		}
		bidangs = append(bidangs, bidang)
	}

	if len(bidangs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada bidang ditemukan"})
		return
	}

	c.JSON(http.StatusOK, bidangs)
}

// Ambil bidang berdasarkan ID
func GetBidangByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid", "detail": err.Error()})
		return
	}

	// Periksa apakah collection sudah diinisialisasi
	if config.BidangCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bidang collection not initialized"})
		return
	}

	var bidang models.Bidang
	err = config.BidangCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&bidang)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bidang tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil bidang", "detail": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, bidang)
}

// Update bidang
func UpdateBidang(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid", "detail": err.Error()})
		return
	}

	// Periksa apakah collection sudah diinisialisasi
	if config.BidangCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bidang collection not initialized"})
		return
	}

	var input models.Bidang
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca input", "detail": err.Error()})
		return
	}

	// Validasi status
	if input.Status != 0 && input.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus 0 (Non-Aktif) atau 1 (Aktif)"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"nama":      input.Nama,
			"status":    input.Status,
			"updatedAt": time.Now(),
		},
	}

	result, err := config.BidangCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui bidang", "detail": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bidang tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bidang berhasil diperbarui"})
}

// Hapus bidang
func DeleteBidang(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid", "detail": err.Error()})
		return
	}

	// Periksa apakah collection sudah diinisialisasi
	if config.BidangCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bidang collection not initialized"})
		return
	}

	res, err := config.BidangCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus bidang", "detail": err.Error()})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bidang tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bidang berhasil dihapus"})
}
