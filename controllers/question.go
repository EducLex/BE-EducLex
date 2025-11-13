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

var questionCollection = config.QuestionCollection

// ✅ POST: Tambah pertanyaan (user)
func CreateQuestion(c *gin.Context) {
	var input models.Question

	// Ambil JSON dari body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tambahkan tanggal dan status default
	input.Tanggal = time.Now()
	input.Status = "Belum Dijawab"

	collection := config.QuestionCollection
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi database belum siap"})
		return
	}

	// Insert ke MongoDB
	result, err := collection.InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ambil ID dari hasil insert
	insertedID := result.InsertedID.(primitive.ObjectID)
	input.ID = insertedID

	c.JSON(http.StatusOK, gin.H{
		"message": "Pertanyaan berhasil ditambahkan",
		"data":    input,
	})
}

// ✅ GET: Semua pertanyaan
func GetQuestions(c *gin.Context) {
	collection := config.QuestionCollection
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi database belum siap"})
		return
	}

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer cursor.Close(context.Background())

	var questions []models.Question
	if err := cursor.All(context.Background(), &questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// ✅ PUT: Update jawaban oleh Jaksa
func UpdateQuestion(c *gin.Context) {
	// Ambil koleksi dari config setiap kali fungsi dipanggil
	collection := config.QuestionCollection
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi database belum siap"})
		return
	}

	// Ambil parameter ID
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Ambil data dari body JSON
	var body struct {
		Jawaban string `json:"jawaban"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// Update field jawaban & status
	update := bson.M{
		"$set": bson.M{
			"jawaban": body.Jawaban,
			"status":  "Sudah Dijawab",
		},
	}

	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jawaban"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertanyaan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jawaban berhasil disimpan"})
}

// ✅ DELETE: Hapus pertanyaan
func DeleteQuestion(c *gin.Context) {
	// Pastikan koneksi database aktif
	collection := config.QuestionCollection
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi database belum siap"})
		return
	}

	// Ambil ID dari parameter URL
	idParam := c.Param("id")
	if len(idParam) != 24 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid (harus 24 karakter hex)"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengonversi ID ke ObjectID: " + err.Error()})
		return
	}

	// Hapus dokumen berdasarkan _id
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pertanyaan: " + err.Error()})
		return
	}

	// Jika tidak ada data yang dihapus
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertanyaan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pertanyaan berhasil dihapus"})
}

// ✅ POST: Tambah diskusi lanjutan
func TambahDiskusi(c *gin.Context) {
	collection := config.QuestionCollection
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi database belum siap"})
		return
	}

	// Ambil parameter ID dan validasi
	idParam := c.Param("id")
	if len(idParam) != 24 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak valid"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengonversi ID ke ObjectID: " + err.Error()})
		return
	}

	// Parsing JSON ke struct Diskusi
	var diskusi models.Diskusi
	if err := c.ShouldBindJSON(&diskusi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	diskusi.Tanggal = time.Now()

	// Siapkan update dasar
	update := bson.M{
		"$push": bson.M{
			"diskusi": diskusi,
		},
	}

	// Jika pengirim adalah Jaksa → ubah status jadi "Sudah Dijawab"
	if diskusi.Pengirim == "Jaksa" {
		update["$set"] = bson.M{
			"status": "Sudah Dijawab",
		}
	}

	// Jalankan update ke MongoDB
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan diskusi: " + err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertanyaan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diskusi berhasil ditambahkan"})
}

// Ambil semua diskusi dari satu pertanyaan berdasarkan ID
func GetDiskusiByQuestionID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	collection := config.QuestionCollection
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi database belum siap"})
		return
	}

	var question models.Question
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&question)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertanyaan tidak ditemukan"})
		return
	}

	// Ambil hanya bagian diskusinya
	c.JSON(http.StatusOK, gin.H{
		"id":      question.ID.Hex(),
		"diskusi": question.Diskusi,
	})
}
