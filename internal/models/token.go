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
	SECRET_KEY string `envconfig:"SECRET_KEY" required:"true"`
}
