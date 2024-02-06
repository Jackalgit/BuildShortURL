package models

import "github.com/google/uuid"

type UserDeleteURL struct {
	UserID   uuid.UUID
	ShortURL string
}

var InputChUserURL = make(chan UserDeleteURL)
