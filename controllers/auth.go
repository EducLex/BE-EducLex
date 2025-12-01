package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/middleware"
	"github.com/EducLex/BE-EducLex/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	}

	// Bind input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi email sudah terdaftar
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, _ := config.UserCollection.CountDocuments(ctx, bson.M{"email": input.Email})
	if count > 0 {
		log.Printf("Email already registered: %v", input.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// Validasi username sudah terdaftar
	count, _ = config.UserCollection.CountDocuments(ctx, bson.M{"username": input.Username})
	if count > 0 {
		log.Printf("Username already registered: %v", input.Username)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar"})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}

	// Buat user baru
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	// Simpan user baru ke database
	_, err = config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user"})
		return
	}

	// Generate OTP untuk verifikasi email
	otp := generateOTP()

	// Kirim OTP ke email pengguna
	message := fmt.Sprintf("Kode OTP Anda untuk verifikasi email: %s", otp)
	err = sendEmail(input.Email, "Verifikasi Email", message)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		c.JSON(500, gin.H{"error": "Gagal mengirim OTP ke email"})
		return
	}

	// Simpan OTP dan waktu kedaluwarsa di database
	_, err = config.UserCollection.UpdateOne(
		ctx,
		bson.M{"email": input.Email},
		bson.M{
			"$set": bson.M{
				"email_verification_otp":    otp,
				"email_verification_expiry": time.Now().Add(10 * time.Minute).Unix(),
			},
		},
	)
	if err != nil {
		log.Printf("Error updating OTP in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan OTP verifikasi"})
		return
	}

	// Generate JWT token
	token, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{
		"message": "Registrasi sukses, cek email Anda untuk verifikasi",
		"token":   token,
	})

}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := config.UserCollection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{"message": "login success", "token": token})
}

// Register Admin
func RegisterAdmin(c *gin.Context) {
	var input struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, _ := config.UserCollection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"email": input.Email},
			{"username": input.Username},
		},
	})
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Email already exists"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "admin",
	}
	_, err := config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}

	token, _ := middleware.GenerateJWT(user.ID.Hex(), user.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{"message": "admin register success", "token": token})
}

// Logout
func Logout(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")
	if tokenHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header tidak ditemukan"})
		return
	}

	token := tokenHeader
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blacklist := models.TokenBlacklist{
		Token:     token,
		ExpiredAt: time.Now().Add(24 * time.Hour), // token dianggap kadaluarsa 1 hari
	}

	_, err := config.TokenBlacklistCollection.InsertOne(ctx, blacklist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan token ke blacklist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

func isValidEmail(email string) bool {
	// Regex untuk validasi email
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
