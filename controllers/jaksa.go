package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
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
	// Pastikan hanya admin yang bisa menambahkan Jaksa
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa menambahkan Jaksa"})
		return
	}

	var body models.Jaksa
	// Ambil JSON dari body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi confirm_password
	if body.Password != body.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password dan Confirm Password tidak cocok"})
		return
	}

	// Hash password sebelum disimpan ke database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}
	body.Password = string(hashedPassword)

	// Validasi jika email atau username sudah terdaftar
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Cek apakah email atau username sudah ada
	count, _ := config.JaksaCollection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"email": body.Email},
			{"username": body.Username},
		},
	})

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username atau Email sudah terdaftar"})
		return
	}

	// Generate OTP untuk verifikasi email
	emailVerificationOTP := generateOTP()
	emailVerificationExpiry := time.Now().Add(10 * time.Minute).Unix()

	// Insert Jaksa baru ke MongoDB (koleksi Jaksa)
	result, err := config.JaksaCollection.InsertOne(context.Background(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan Jaksa"})
		return
	}

	// Mengambil ID yang baru dimasukkan untuk Jaksa
	insertedID := result.InsertedID.(primitive.ObjectID)
	body.ID = insertedID

	// Update Jaksa dengan OTP verifikasi email
	_, err = config.JaksaCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": body.ID},
		bson.M{
			"$set": bson.M{
				"email_verification_otp":    emailVerificationOTP,
				"email_verification_expiry": emailVerificationExpiry,
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan OTP verifikasi email"})
		return
	}

	// Kirim email verifikasi
	err = sendVerificationEmail(body.Email, emailVerificationOTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim email verifikasi"})
		return
	}

	// Tambahkan Jaksa ke koleksi User (buat entry di User)
	user := models.User{
		ID:       body.ID,
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password, 
		Role:     "jaksa",       
		EmailVerificationOTP:    emailVerificationOTP,
		EmailVerificationExpiry: emailVerificationExpiry,
	}

	_, err = config.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan Jaksa ke koleksi User"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Jaksa berhasil ditambahkan ke koleksi Jaksa dan User, email verifikasi telah dikirim",
		"data":    body,
	})
}

func sendVerificationEmail(to, otp string) error {
	// Email pengirim dan App Password
	from := "dewidesember20@gmail.com" // Ganti dengan email Gmail kamu
	pass := "pezf jucw gssc mmar"      // Ganti dengan App Password yang kamu buat

	// Mengatur SMTP server Gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Format pesan email
	message := fmt.Sprintf("Kode OTP Anda untuk verifikasi email: %s", otp)
	msg := fmt.Sprintf("Subject: Verifikasi Email\r\n\r\n%s", message)

	// Mengonfigurasi otentikasi SMTP
	auth := smtp.PlainAuth("", from, pass, smtpHost)

	// Kirim email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	log.Println("Email sent successfully!")
	return nil
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
