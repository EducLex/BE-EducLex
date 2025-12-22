package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username                string             `bson:"username" json:"username"`
	Email                   string             `bson:"email" json:"email"`
	Password                string             `bson:"password,omitempty" json:"-"`
	GoogleID                string             `bson:"google_id,omitempty" json:"google_id"`
	Role                    string             `bson:"role,omitempty" json:"role"`
	Token                   string             `bson:"token,omitempty" json:"token"`
	JaksaID                 primitive.ObjectID `bson:"jaksa_id,omitempty"`
	ResetOtp                string             `bson:"reset_otp,omitempty" json:"reset_otp,omitempty"`
	ResetOtpExpiry          int64              `bson:"reset_otp_expiry,omitempty" json:"reset_otp_expiry,omitempty"`
	EmailVerified           bool               `bson:"email_verified" json:"email_verified"`
	EmailVerificationOTP    string             `bson:"email_verification_otp,omitempty" json:"email_verification_otp,omitempty"`
	EmailVerificationExpiry int64              `bson:"email_verification_expiry,omitempty" json:"email_verification_expiry,omitempty"`
}

type DashboardData struct {
	TotalArtikel   int `json:"totalArtikel"`
	TotalTanya     int `json:"totalTanya"`
	TotalTulisan   int `json:"totalTulisan"`
	TotalPeraturan int `json:"totalPeraturan"`
}

type Diskusi struct {
	Pengirim string    `json:"pengirim" bson:"pengirim"`
	Pesan    string    `json:"pesan" bson:"pesan"`
	Tanggal  time.Time `json:"tanggal" bson:"tanggal"`
}

type Question struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Nama       string             `json:"nama" bson:"nama"`
	Email      string             `json:"email,omitempty" bson:"email,omitempty"`
	Kategori   string             `json:"kategori" bson:"kategori"`
	Pertanyaan string             `json:"pertanyaan" bson:"pertanyaan"`
	Jawaban    string             `json:"jawaban,omitempty" bson:"jawaban,omitempty"`
	Status     string             `json:"status" bson:"status"`
	Tipe       string             `json:"tipe,omitempty" bson:"tipe,omitempty"`
	Tanggal    time.Time          `json:"tanggal" bson:"tanggal"`
	Diskusi    []Diskusi          `json:"diskusi,omitempty" bson:"diskusi,omitempty"`
	BidangID   primitive.ObjectID `json:"bidang_id" bson:"bidang_id"`
	BidangNama string             `json:"bidang_nama" bson:"bidang_nama"`
}

type Article struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Judul      string             `bson:"judul" json:"judul"`
	Isi        string             `bson:"isi" json:"isi"`
	Penulis    string             `bson:"penulis" json:"penulis"`
	Gambar     string             `bson:"gambar,omitempty" json:"gambar,omitempty"`
	Dokumen    string             `bson:"dokumen,omitempty" json:"dokumen,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	CategoryID primitive.ObjectID `bson:"categoryId" json:"categoryId"`
}

type Tulisan struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Penulis    string             `bson:"penulis" json:"penulis"`
	Judul      string             `bson:"judul" json:"judul"`
	Isi        string             `bson:"isi" json:"isi"`
	File       string             `bson:"file,omitempty" json:"file,omitempty"`
	Gambar     string             `bson:"gambar,omitempty" json:"gambar,omitempty"`
	Dokumen    string             `bson:"dokumen,omitempty" json:"dokumen,omitempty"`
	BidangID   primitive.ObjectID `bson:"bidang_id" json:"bidang_id"`
	BidangNama string             `bson:"bidang_nama" json:"bidang_nama"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Peraturan struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Judul     string    `bson:"judul" json:"judul"`
	Isi       string    `bson:"isi" json:"isi"`
	Kategori  string    `bson:"kategori" json:"kategori"`
	Gambar    string    `bson:"gambar,omitempty" json:"gambar,omitempty"`
	Dokumen   string    `bson:"dokumen,omitempty" json:"dokumen,omitempty"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type TokenBlacklist struct {
	Token     string    `bson:"token"`
	ExpiredAt time.Time `bson:"expired_at"`
}

type Jaksa struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nama                    string             `bson:"nama" json:"nama"`
	Username                string             `bson:"username" json:"username"`
	Email                   string             `bson:"email" json:"email"`
	NIP                     string             `bson:"nip" json:"nip"`
	UserID                  primitive.ObjectID `bson:"user_id" json:"user_id"`
	Foto                    string             `bson:"foto,omitempty" json:"foto,omitempty"`
	Password                string             `json:"password,omitempty" bson:"password,omitempty"`
	ConfirmPassword         string             `bson:"-" json:"confirm_password,omitempty"`
	ResetOtp                string             `bson:"reset_otp,omitempty" json:"reset_otp,omitempty"`
	ResetOtpExpiry          int64              `bson:"reset_otp_expiry,omitempty" json:"reset_otp_expiry,omitempty"`
	BidangID                primitive.ObjectID `json:"bidang_id" bson:"bidang_id"`
	BidangNama              string             `json:"bidang_nama" bson:"bidang_nama"`
	EmailVerificationOTP    string             `bson:"email_verification_otp,omitempty" json:"email_verification_otp,omitempty"`
	EmailVerificationExpiry int64              `bson:"email_verification_expiry,omitempty" json:"email_verification_expiry,omitempty"`
	EmailVerified           bool               `bson:"email_verified" json:"email_verified"`
}

type UpdateJaksaRequest struct {
	Nama     string `json:"nama"`
	NIP      string `json:"nip"`
	Email    string `json:"email"`
	BidangID string `json:"bidang_id"`
}

type Category struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Subkategori string             `bson:"subkategori" json:"subkategori"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Bidang struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Nama   string             `json:"nama" bson:"nama"`
	Status int                `json:"status" bson:"status"`
}
