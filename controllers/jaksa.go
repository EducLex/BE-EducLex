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
	"golang.org/x/crypto/bcrypt"
)

// =========================
// Create Jaksa
// =========================
func CreateJaksa(c *gin.Context) {
	var body models.Jaksa

	// Ambil JSON dari body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Hash password sebelum disimpan ke database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}
	body.Password = string(hashedPassword)

	// Insert Jaksa baru ke MongoDB (koleksi Jaksa)
	result, err := config.JaksaCollection.InsertOne(context.Background(), body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal menambahkan Jaksa"})
		return
	}

	// Mengambil ID yang baru dimasukkan untuk Jaksa
	insertedID := result.InsertedID.(primitive.ObjectID)
	body.ID = insertedID

	// Tambahkan Jaksa ke koleksi User (buat entry di User)
	user := models.User{
		ID:       body.ID,
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,  
		Role:     "jaksa",       
	}

	_, err = config.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal menambahkan Jaksa ke koleksi User"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Jaksa berhasil ditambahkan ke koleksi Jaksa dan User",
		"data":    body,
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

	// Ambil userId dari middleware autentikasi
	userId, _ := c.Get("userId")
	if userId == nil {
		c.JSON(400, gin.H{"error": "User tidak terautentikasi"})
		return
	}

	// Set userId ke dalam objek Jaksa
	body.UserID = userId.(primitive.ObjectID)

	update := bson.M{
		"$set": bson.M{
			"nama":    body.Nama,
			"nip":     body.NIP,
			"jabatan": body.Jabatan,
			"foto":    body.Foto,
			"user_id": body.UserID, // Update userId juga
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
