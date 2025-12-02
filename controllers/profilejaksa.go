package controllers

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
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
	// Ambil ID dari parameter URL
	id := c.Param("id")

	// Konversi ID menjadi ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Jika ID tidak valid, kembalikan error
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	// Buat variabel untuk menyimpan data Jaksa
	var data models.Jaksa

	// Gunakan context dengan timeout untuk efisiensi query
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cari data Jaksa berdasarkan ID
	err = config.JaksaCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&data)
	if err != nil {
		// Jika data tidak ditemukan
		if err.Error() == "mongo: no documents in result" {
			c.JSON(404, gin.H{"error": "Data tidak ditemukan"})
			return
		}
		// Jika ada error lain
		c.JSON(500, gin.H{"error": "Gagal mengambil data Jaksa"})
		return
	}

	// Kembalikan response dengan data Jaksa
	c.JSON(200, gin.H{"data": data})
}

// =========================
// Update Jaksa Profile
// =========================
func UpdateJaksaProfile(c *gin.Context) {
	id := c.Param("id")

	// Convert id to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID tidak valid"})
		return
	}

	// Struct untuk input data baru
	var body struct {
		Nama       string `json:"nama"`
		NIP        string `json:"nip"`
		Email      string `json:"email"`
		Foto       string `json:"foto"`
		BidangID   string `json:"bidang_id"`
		BidangNama string `json:"bidang_nama"`
	}

	// Bind input dari JSON body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	// Memastikan BidangID dan BidangNama ada
	if body.BidangID == "" || body.BidangNama == "" {
		c.JSON(400, gin.H{"error": "Bidang ID dan Nama harus diisi"})
		return
	}

	// Siapkan data update
	update := bson.M{
		"$set": bson.M{
			"nama":        body.Nama,
			"nip":         body.NIP,
			"email":       body.Email,
			"foto":        body.Foto,
			"bidang_id":   body.BidangID,   // Update BidangID
			"bidang_nama": body.BidangNama, // Update BidangNama
		},
	}

	// Cek dan update database
	_, err = config.JaksaCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		log.Printf("Error updating Jaksa: %v", err)
		c.JSON(500, gin.H{"error": "Update gagal", "detail": err.Error()})
		return
	}

	// Jika berhasil update
	c.JSON(200, gin.H{"message": "Profil berhasil diperbarui"})
}

// Fungsi untuk generate OTP secara acak
func generateOTP() string {
	// Membuat OTP acak 6 digit
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}

// Fungsi untuk verifikasi email (untuk User dan Jaksa)
func VerifyEmail(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	// Bind request body ke struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Data tidak valid"})
		return
	}

	// Cari pengguna berdasarkan email di koleksi User
	var user models.User
	err := config.UserCollection.FindOne(context.Background(), bson.M{"email": body.Email}).Decode(&user)
	if err == nil {
		// Periksa apakah OTP sesuai
		if user.EmailVerificationOTP != body.OTP {
			c.JSON(400, gin.H{"error": "OTP tidak valid"})
			return
		}

		// Cek apakah OTP sudah kedaluwarsa
		if time.Now().Unix() > user.EmailVerificationExpiry {
			c.JSON(400, gin.H{"error": "OTP sudah kedaluwarsa"})
			return
		}

		// Update status email terverifikasi untuk User
		_, err = config.UserCollection.UpdateOne(
			context.Background(),
			bson.M{"email": body.Email},
			bson.M{
				"$set": bson.M{"email_verified": true},
				"$unset": bson.M{
					"email_verification_otp":    "",
					"email_verification_expiry": "",
				},
			},
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "Gagal memperbarui status verifikasi email"})
			return
		}

		c.JSON(200, gin.H{"message": "Email berhasil diverifikasi (User)"})
		return
	}

	// Jika tidak ditemukan di User, coba cari di koleksi Jaksa
	var jaksa models.Jaksa
	err = config.JaksaCollection.FindOne(context.Background(), bson.M{"email": body.Email}).Decode(&jaksa)
	if err != nil {
		c.JSON(404, gin.H{"error": "Email tidak terdaftar sebagai User atau Jaksa"})
		return
	}

	// Periksa apakah OTP sesuai dengan Jaksa
	if jaksa.EmailVerificationOTP != body.OTP {
		c.JSON(400, gin.H{"error": "OTP tidak valid"})
		return
	}

	// Cek apakah OTP sudah kedaluwarsa untuk Jaksa
	if time.Now().Unix() > jaksa.EmailVerificationExpiry {
		c.JSON(400, gin.H{"error": "OTP sudah kedaluwarsa"})
		return
	}

	// Update status email terverifikasi untuk Jaksa
	_, err = config.JaksaCollection.UpdateOne(
		context.Background(),
		bson.M{"email": body.Email},
		bson.M{
			"$set": bson.M{"email_verified": true},
			"$unset": bson.M{
				"email_verification_otp":    "",
				"email_verification_expiry": "",
			},
		},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal memperbarui status verifikasi email Jaksa"})
		return
	}

	c.JSON(200, gin.H{"message": "Email berhasil diverifikasi (Jaksa)"})
}

// Fungsi untuk mengirim email dengan SMTP Gmail
func sendEmail(to, subject, message string) error {
	from := "dewidesember20@gmail.com"
	pass := "pezf jucw gssc mmar"

	// Mengatur SMTP server Gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Format pesan email
	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, message)

	// Membuat koneksi ke server SMTP
	serverAddr := smtpHost + ":" + smtpPort
	conn, err := net.DialTimeout("tcp", serverAddr, 30*time.Second)
	if err != nil {
		log.Println("Error dialing server:", err)
		return fmt.Errorf("gagal menghubungi server SMTP: %w", err)
	}
	defer conn.Close()

	// Membuat konfigurasi TLS untuk koneksi
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	// Menggunakan koneksi TLS
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Println("Error creating SMTP client:", err)
		return fmt.Errorf("gagal membuat client SMTP: %w", err)
	}

	// Memulai sesi TLS
	if err := client.StartTLS(tlsConfig); err != nil {
		log.Println("Error starting TLS:", err)
		return fmt.Errorf("gagal memulai sesi TLS: %w", err)
	}

	// Autentikasi dengan server SMTP
	auth := smtp.PlainAuth("", from, pass, smtpHost)
	if err := client.Auth(auth); err != nil {
		log.Println("Error authenticating with SMTP:", err)
		return fmt.Errorf("gagal otentikasi ke SMTP: %w", err)
	}

	// Menentukan pengirim dan penerima
	if err := client.Mail(from); err != nil {
		log.Println("Error setting sender:", err)
		return fmt.Errorf("gagal menetapkan pengirim: %w", err)
	}

	// Menambahkan penerima email
	if err := client.Rcpt(to); err != nil {
		log.Println("Error setting recipient:", err)
		return fmt.Errorf("gagal menetapkan penerima: %w", err)
	}

	// Menulis email ke dalam koneksi
	writer, err := client.Data()
	if err != nil {
		log.Println("Error opening writer:", err)
		return fmt.Errorf("gagal membuka penulis: %w", err)
	}

	_, err = writer.Write([]byte(msg))
	if err != nil {
		log.Println("Error writing email:", err)
		return fmt.Errorf("gagal menulis pesan email: %w", err)
	}

	// Mengakhiri pengiriman email
	err = writer.Close()
	if err != nil {
		log.Println("Error closing writer:", err)
		return fmt.Errorf("gagal menutup penulis: %w", err)
	}

	// Mengakhiri koneksi SMTP
	client.Quit()

	log.Println("Email berhasil dikirim!")
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

	// Cek apakah email sudah terverifikasi
	if user.EmailVerificationExpiry > time.Now().Unix() {
		// Jika email belum terverifikasi, kembalikan pesan error
		c.JSON(400, gin.H{"error": "Email belum terverifikasi"})
		return
	}

	// Generate OTP untuk reset password
	otp := generateOTP()

	// Kirim OTP ke email pengguna
	message := fmt.Sprintf("Kode OTP Anda untuk mereset password: %s", otp)
	err = sendEmail(body.Email, "Reset Password", message)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal mengirim OTP ke email"})
		return
	}

	// Set OTP dan waktu kedaluwarsa di database
	_, err = config.UserCollection.UpdateOne(
		context.Background(),
		bson.M{"email": body.Email},
		bson.M{
			"$set": bson.M{
				"reset_otp":        otp,
				"reset_otp_expiry": time.Now().Add(10 * time.Minute).Unix(),
			},
		},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal menyimpan OTP reset password"})
		return
	}

	c.JSON(200, gin.H{"message": "OTP untuk reset password sudah dikirim ke email"})
}

// Fungsi untuk mereset password pengguna di Jaksa
func ResetPassword(c *gin.Context) {
	var body struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}

	// Bind request body ke struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Data tidak valid"})
		return
	}

	// Cari pengguna berdasarkan email dan OTP
	var user models.User
	err := config.UserCollection.FindOne(context.Background(), bson.M{
		"email":     body.Email,
		"reset_otp": body.OTP,
	}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"error": "OTP tidak valid"})
		return
	}

	// Cek apakah OTP sudah kedaluwarsa
	if time.Now().Unix() > user.ResetOtpExpiry {
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
	_, err = config.UserCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
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

	c.JSON(200, gin.H{"message": "Password berhasil direset"})
}
