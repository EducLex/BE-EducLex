package controllers

import (
	"context"
	"net/http"
	"time"
	"os"
	"fmt"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ✅ Tambah peraturan baru (Admin)
func CreatePeraturan(c *gin.Context) {
	// Membuat folder uploads jika belum ada
	os.MkdirAll("uploads", os.ModePerm)

	var input models.Peraturan
	// Mengambil data dari Form Data (bukan JSON)
	input.Judul = c.PostForm("judul")
	input.Isi = c.PostForm("isi")

	// Mengambil kategori dari Form Data
	kategori := c.PostForm("kategori")
	if kategori != "internal" && kategori != "eksternal" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori harus 'internal' atau 'eksternal'"})
		return
	}

	// Validasi kategori dan subkategori
	if kategori == "eksternal" {
		// Eksternal: Peraturan UUD, Peraturan Presiden, Perpres, Keppres
		subKategori := c.PostForm("subkategori")
		if subKategori != "Peraturan UUD" && subKategori != "Peraturan Presiden" && subKategori != "Perpres" && subKategori != "Keppres" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subkategori eksternal tidak valid"})
			return
		}
		input.Kategori = fmt.Sprintf("%s: %s", kategori, subKategori)
	} else {
		// Internal: Pembinaan, Intelijen, Pidana Umum, dll
		subKategori := c.PostForm("subkategori")
		internalCategories := []string{"Pembinaan", "Intelijen", "Pidana Umum", "Pidana Khusus", "Perdata dan Tata Usaha Negara", "Pidana Militer", "Asisten Pengawasan"}
		valid := false
		for _, cat := range internalCategories {
			if subKategori == cat {
				valid = true
				break
			}
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Subkategori internal tidak valid"})
			return
		}
		input.Kategori = fmt.Sprintf("%s: %s", kategori, subKategori)
	}

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

	// Insert peraturan ke MongoDB
	collection := config.PeraturanCollection
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	res, err := collection.InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Gagal menyimpan peraturan",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Peraturan berhasil ditambahkan",
		"id":      res.InsertedID,
	})
}

// ✅ Ambil semua peraturan (User & Admin)
func GetPeraturan(c *gin.Context) {
	// 🔥 PAKSA CORS HEADER DI SINI
	c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5501")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, Accept")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	// (untuk jaga-jaga kalau ada OPTIONS)
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(204)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.PeraturanCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer cursor.Close(ctx)

	var results []models.Peraturan
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal decode data"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ✅ Ambil peraturan berdasarkan ID (User & Admin)
func GetPeraturanByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var peraturan models.Peraturan
	err = config.PeraturanCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&peraturan)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peraturan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, peraturan)
}

// ✅ Update peraturan (Admin)
func UpdatePeraturan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input models.Peraturan
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"judul":     input.Judul,
			"isi":       input.Isi,
			"kategori":  input.Kategori,
			"updatedAt": time.Now(),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = config.PeraturanCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peraturan berhasil diperbarui"})
}

// ✅ Hapus peraturan (Admin)
func DeletePeraturan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = config.PeraturanCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peraturan berhasil dihapus"})
}
