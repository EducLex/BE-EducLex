package models

import "time"

type User struct {
	Email     string    `bson:"email" json:"email"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password,omitempty" json:"password,omitempty"`
	Provider  string    `bson:"provider" json:"provider"` // "manual" / "google"
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
