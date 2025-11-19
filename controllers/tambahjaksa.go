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

// Tambah Jaksa
func CreateJaksa(c *gin.Context) {
	jaksaCollection := config.JaksaCollection
	if jaksaCollection == nil {
		c.JSON(500, gin.H{"error": "JaksaCollection masih nil"})
		return
	}

	var input models.Jaksa
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := jaksaCollection.InsertOne(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data jaksa: " + err.Error()})
		return
	}

	input.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Data Jaksa berhasil ditambahkan",
		"data":    input,
	})
}

// Get All Jaksa
func GetAllJaksa(c *gin.Context) {
	jaksaCollection := config.JaksaCollection
	if jaksaCollection == nil {
		c.JSON(500, gin.H{"error": "JaksaCollection masih nil"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := jaksaCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data jaksa"})
		return
	}
	defer cursor.Close(ctx)

	var list []models.Jaksa
	if err := cursor.All(ctx, &list); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data jaksa"})
		return
	}

	c.JSON(200, gin.H{"data": list})
}

// Update Jaksa
func UpdateJaksa(c *gin.Context) {
	jaksaCollection := config.JaksaCollection
	if jaksaCollection == nil {
		c.JSON(500, gin.H{"error": "JaksaCollection masih nil"})
		return
	}

	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	var body models.Jaksa
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Data tidak valid"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"nama":    body.Nama,
			"nip":     body.NIP,
			"jabatan": body.Jabatan,
			"email":   body.Email,
			"foto":    body.Foto,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := jaksaCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal update jaksa"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "Jaksa tidak ditemukan"})
		return
	}

	c.JSON(200, gin.H{"message": "Data Jaksa berhasil diperbarui"})
}

// Delete Jaksa
func DeleteJaksa(c *gin.Context) {
	jaksaCollection := config.JaksaCollection
	if jaksaCollection == nil {
		c.JSON(500, gin.H{"error": "JaksaCollection masih nil"})
		return
	}

	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := jaksaCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal menghapus data jaksa"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "Jaksa tidak ditemukan"})
		return
	}

	c.JSON(200, gin.H{"message": "Data Jaksa berhasil dihapus"})
}
