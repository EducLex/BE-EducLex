package controllers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
// Get Jaksa Profile
// =========================
func GetJaksaProfile(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	var data models.Jaksa
	err = config.JaksaCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&data)
	if err != nil {
		c.JSON(404, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	data.Password = ""
	c.JSON(200, gin.H{"data": data})
}

// =========================
// Update Jaksa Profile
// =========================
func UpdateJaksaProfile(c *gin.Context) {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	var body struct {
		Nama  string `json:"nama"`
		NIP   string `json:"nip"`
		Email string `json:"email"`
		Foto  string `json:"foto"`
	}

	if c.ShouldBindJSON(&body) != nil {
		c.JSON(400, gin.H{"error": "Data tidak valid"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"nama":  body.Nama,
			"nip":   body.NIP,
			"email": body.Email,
			"foto":  body.Foto,
		},
	}

	_, err = config.JaksaCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(500, gin.H{"error": "Update gagal"})
		return
	}

	c.JSON(200, gin.H{"message": "Profil berhasil diperbarui"})
}

// Fungsi untuk generate OTP secara acak
func generateOTP() string {
	// Membuat OTP acak 6 digit
	otp := fmt.Sprintf("%06d", rand.Intn(1000000)) // Generate angka acak antara 000000 hingga 999999
	return otp
}

// Fungsi untuk mengirim email dengan SMTP Gmail
func sendEmail(to, subject, message string) error {
	// Email pengirim dan App Password
	from := "dewidesember20@gmail.com" // Ganti dengan email Gmail kamu
	pass := "pezf jucw gssc mmar"      // Ganti dengan App Password yang kamu buat

	// Mengatur SMTP server Gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Format pesan email
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, message)

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

// Fungsi ForgotPassword
func ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	// Bind request body ke struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Email tidak valid"})
		return
	}

	// Cari pengguna berdasarkan email di UserCollection
	var user models.User
	err := config.UserCollection.FindOne(context.Background(), bson.M{"email": body.Email}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"error": "Email tidak ditemukan"})
		return
	}

	// Generate OTP secara dinamis
	otp := generateOTP()

	// Set OTP expiry to 10 minutes
	expiry := time.Now().Add(10 * time.Minute).Unix()

	// Simpan OTP dan waktu kedaluwarsa di database JaksaCollection
	_, err = config.JaksaCollection.UpdateOne(
		context.Background(),
		bson.M{"email": body.Email}, // Menyimpan OTP di collection Jaksa berdasarkan email
		bson.M{
			"$set": bson.M{
				"reset_otp":        otp,
				"reset_otp_expiry": expiry,
			},
		},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengupdate OTP"})
		return
	}

	// Kirim OTP ke email pengguna
	message := fmt.Sprintf("Kode OTP Anda untuk mereset password: %s", otp)
	err = sendEmail(body.Email, "Reset Password", message)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengirim OTP ke email"})
		return
	}

	// Response sukses
	c.JSON(200, gin.H{"message": "OTP untuk reset password sudah dikirim ke email"})
}

// Fungsi untuk mereset password pengguna di Jaksa
func ResetPassword(c *gin.Context) {
	var body struct {
		Email       string `json:"email"` // Email pengguna untuk verifikasi OTP
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}

	// Bind request body ke struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Data tidak valid"})
		return
	}

	// Cari Jaksa berdasarkan email dan OTP di JaksaCollection
	var jaksa models.Jaksa
	err := config.JaksaCollection.FindOne(context.Background(), bson.M{
		"email":    body.Email,
		"reset_otp": body.OTP,
	}).Decode(&jaksa)
	if err != nil {
		c.JSON(404, gin.H{"error": "OTP tidak valid"})
		return
	}

	// Log OTP yang disimpan di database dan yang dimasukkan pengguna
	fmt.Println("OTP yang disimpan di database:", jaksa.ResetOtp)
	fmt.Println("OTP yang dimasukkan pengguna:", body.OTP)

	// Cek apakah OTP sudah kedaluwarsa
	if time.Now().Unix() > jaksa.ResetOtpExpiry {
		c.JSON(400, gin.H{"error": "OTP sudah kedaluwarsa"})
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal meng-hash password"})
		return
	}

	// Update password dan hapus OTP
	_, err = config.JaksaCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": jaksa.ID},
		bson.M{
			"$set": bson.M{"password": string(hashedPassword)},
			"$unset": bson.M{
				"reset_otp":        "",
				"reset_otp_expiry": "",
			},
		},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mereset password"})
		return
	}

	// Respons sukses
	c.JSON(200, gin.H{"message": "Password berhasil direset"})
}
