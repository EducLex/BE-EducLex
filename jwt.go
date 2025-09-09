package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("SECRET_KEY")

func GenerateJWT(email, name string) (string, error) {
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
