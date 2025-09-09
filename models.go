package main

type User struct {
	Email string `bson:"email" json:"email"`
	Name  string `bson:"name"  json:"name"`
}
