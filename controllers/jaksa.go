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

	// Prepare data from form inputs
	var body models.Jaksa
	body.Username = c.PostForm("username")
	body.Nama = c.PostForm("nama")
	body.Email = c.PostForm("email")
	body.NIP = c.PostForm("nip")

	// Validate confirm_password
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")
	if password != confirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password dan Confirm Password tidak cocok"})
		return
	}

	// Hash password before storing in DB
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}
	body.Password = string(hashedPassword)

	// Validate if email or username already exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	// Generate OTP for email verification
	emailVerificationOTP := generateOTP()
	emailVerificationExpiry := time.Now().Add(10 * time.Minute).Unix()

	// Handling BidangID and getting BidangNama
	bidangID := c.PostForm("bidang_id")
	if bidangID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bidang ID tidak boleh kosong"})
		return
	}

	// Convert bidang_id into ObjectID
	objectID, err := primitive.ObjectIDFromHex(bidangID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bidang ID tidak valid"})
		return
	}
	body.BidangID = objectID

	// Get BidangName based on BidangID
	var bidang models.Bidang
	err = config.BidangCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&bidang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bidang tidak ditemukan"})
		return
	}

	// Automatically set BidangNama based on BidangID
	body.BidangNama = bidang.Nama

	// Insert new Jaksa into MongoDB
	result, err := config.JaksaCollection.InsertOne(context.Background(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan Jaksa"})
		return
	}

	// Get inserted ID for Jaksa
	insertedID := result.InsertedID.(primitive.ObjectID)
	body.ID = insertedID

	// Update Jaksa with email verification OTP
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

	// Send email verification
	err = sendVerificationEmail(body.Email, emailVerificationOTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim email verifikasi"})
		return
	}

	// Add Jaksa to User collection
	user := models.User{
		ID:                      body.ID,
		Username:                body.Username,
		Email:                   body.Email,
		Password:                body.Password,
		Role:                    "jaksa",
		EmailVerificationOTP:    emailVerificationOTP,
		EmailVerificationExpiry: emailVerificationExpiry,
	}

	_, err = config.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan Jaksa ke koleksi User"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message": "Jaksa berhasil ditambahkan ke koleksi Jaksa dan User, email verifikasi telah dikirim",
		"data":    body,
	})
}

func sendVerificationEmail(to, otp string) error {
	// Email pengirim dan App Password
	from := "dewidesember20@gmail.com"
	pass := "pezf jucw gssc mmar"

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

type UpdateJaksaRequest struct {
	Nama     string `json:"nama"`
	NIP      string `json:"nip"`
	Email    string `json:"email"`
	BidangID string `json:"bidang_id"`
}

// Update Jaksa
func UpdateJaksa(c *gin.Context) {
	jaksaCollection := config.JaksaCollection
	userCollection := config.UserCollection

	idParam := c.Param("id")
	jaksaID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID Jaksa tidak valid"})
		return
	}

	// Bind request
	var body UpdateJaksaRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ambil data jaksa lama
	var jaksa models.Jaksa
	if err := jaksaCollection.FindOne(ctx, bson.M{"_id": jaksaID}).Decode(&jaksa); err != nil {
		c.JSON(404, gin.H{"error": "Jaksa tidak ditemukan"})
		return
	}

	updateJaksa := bson.M{}
	updateUser := bson.M{}
	emailDiubah := false

	// Update nama
	if body.Nama != "" {
		updateJaksa["nama"] = body.Nama
		updateUser["username"] = body.Nama // opsional, kalau username = nama
	}

	// Update NIP
	if body.NIP != "" {
		updateJaksa["nip"] = body.NIP
	}

	// Update Email
	if body.Email != "" && body.Email != jaksa.Email {
		emailDiubah = true
		updateJaksa["email"] = body.Email
		updateUser["email"] = body.Email
	}

	// Update Bidang
	if body.BidangID != "" {
		bidangObjID, err := primitive.ObjectIDFromHex(body.BidangID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Bidang ID tidak valid"})
			return
		}

		var bidang models.Bidang
		if err := config.BidangCollection.FindOne(ctx, bson.M{"_id": bidangObjID}).Decode(&bidang); err != nil {
			c.JSON(404, gin.H{"error": "Bidang tidak ditemukan"})
			return
		}

		updateJaksa["bidang_id"] = bidangObjID
		updateJaksa["bidang_nama"] = bidang.Nama
	}

	// Kalau tidak ada data diupdate
	if len(updateJaksa) == 0 {
		c.JSON(400, gin.H{"error": "Tidak ada data yang diupdate"})
		return
	}

	// 🔁 OTP ulang jika email berubah
	if emailDiubah {
		otp := generateOTP()
		expiry := time.Now().Add(10 * time.Minute).Unix()

		updateJaksa["email_verification_otp"] = otp
		updateJaksa["email_verification_expiry"] = expiry
		updateJaksa["email_verified"] = false

		updateUser["email_verification_otp"] = otp
		updateUser["email_verification_expiry"] = expiry
		updateUser["email_verified"] = false

		// Kirim email
		go sendVerificationEmail(body.Email, otp)
	}

	// Update Jaksa
	_, err = jaksaCollection.UpdateOne(
		ctx,
		bson.M{"_id": jaksaID},
		bson.M{"$set": updateJaksa},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal update data Jaksa"})
		return
	}

	// Update User (berdasarkan email lama / jaksa_id)
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"email": jaksa.Email},
		bson.M{"$set": updateUser},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal update data User"})
		return
	}

	c.JSON(200, gin.H{
		"message":                     "Data Jaksa & User berhasil diperbarui",
		"email_verification_required": emailDiubah,
	})
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
