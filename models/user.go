package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Username     string             `bson:"name" json:"name"`
	Password     string             `bson:"password,omitempty" json:"-"`
	GoogleID     string             `bson:"google_id" json:"google_id"`
	RefreshToken string             `bson:"refresh_token" json:"refresh_token"`
}
