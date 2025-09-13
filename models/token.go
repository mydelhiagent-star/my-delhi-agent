package models

import "time"

// TokenData represents a database-agnostic token model
type Token struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
