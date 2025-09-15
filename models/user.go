package models

import (
		"go.mongodb.org/mongo-driver/bson/primitive"
		"time"
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

type Question struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nama       string             `bson:"nama" json:"nama"`
	Email      string             `bson:"email" json:"email"`
	Pertanyaan string             `bson:"pertanyaan" json:"pertanyaan"`
	Jawaban    string             `bson:"jawaban,omitempty" json:"jawaban,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
