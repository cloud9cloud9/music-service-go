package models

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type Token struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Status    string    `json:"status"`
	UserID    int       `json:"user_id"`
}
