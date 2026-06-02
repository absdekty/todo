package model

import "time"

type Token struct {
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	Revoked   bool      `json:"revoked"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
