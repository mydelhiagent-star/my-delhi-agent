package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	ID    string `json:"id"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type Token struct {
	Token string `bson:"token" json:"token"`
	User  string `bson:"user" json:"user"`
}
