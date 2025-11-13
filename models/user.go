package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password,omitempty" json:"-"`
	GoogleID string             `bson:"google_id,omitempty" json:"google_id"`
	Role     string             `bson:"role,omitempty" json:"role"`
	Token    string             `bson:"token,omitempty" json:"token"`
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
	Tipe       string             `json:"tipe,omitempty" bson:"tipe,omitempty"` // "publik" / "internal"
	Tanggal    time.Time          `json:"tanggal" bson:"tanggal"`
	Diskusi    []Diskusi          `json:"diskusi,omitempty" bson:"diskusi,omitempty"`
}

type Article struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Judul     string             `bson:"judul" json:"judul"`
	Isi       string             `bson:"isi" json:"isi"`
	Penulis   string             `bson:"penulis" json:"penulis"`
	Gambar    string             `bson:"gambar,omitempty" json:"gambar,omitempty"`
	Dokumen   string             `bson:"dokumen,omitempty" json:"dokumen,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Tulisan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Penulis   string             `bson:"penulis" json:"penulis"`
	Kategori  string             `bson:"kategori" json:"kategori"`
	Judul     string             `bson:"judul" json:"judul"`
	Isi       string             `bson:"isi" json:"isi"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Peraturan struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Judul     string    `bson:"judul" json:"judul"`
	Isi       string    `bson:"isi" json:"isi"`
	Kategori  string    `bson:"kategori" json:"kategori"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type TokenBlacklist struct {
	Token     string    `bson:"token"`
	ExpiredAt time.Time `bson:"expired_at"`
}

type Jaksa struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nama     string             `bson:"nama" json:"nama"`
	NIP      string             `bson:"nip" json:"nip"`
	Jabatan  string             `bson:"jabatan" json:"jabatan"`
	Email    string             `bson:"email" json:"email"`
	Foto     string             `bson:"foto,omitempty" json:"foto,omitempty"` // optional
}