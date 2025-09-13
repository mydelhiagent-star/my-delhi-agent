package gorm_models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Token     string         `gorm:"size:255;uniqueIndex;not null" json:"token"`
	UserID    string         `gorm:"size:255;not null;index" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"expires_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
