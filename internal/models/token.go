package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

type Secret struct {
	SecretKey string `envconfig:"Secret_key" required:"true"`
}
