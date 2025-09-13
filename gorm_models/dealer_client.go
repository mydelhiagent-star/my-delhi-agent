package gorm_models

import (
	"time"

	"gorm.io/gorm"
)

type DealerClient struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	DealerID   uint           `gorm:"not null;index" json:"dealer_id"`
	PropertyID uint           `gorm:"not null;index" json:"property_id"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Phone      string         `gorm:"size:20;not null" json:"phone"`
	Email      string         `gorm:"size:255" json:"email"`
	Status     string         `gorm:"size:50;default:'active'" json:"status"`
	Note       string         `gorm:"type:text" json:"note"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Dealer   Dealer   `gorm:"foreignKey:DealerID;references:ID" json:"dealer,omitempty"`
	Property Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`
}

// Dealer client status constants
const (
	DealerClientStatusActive   = "active"
	DealerClientStatusInactive = "inactive"
	DealerClientStatusBlocked  = "blocked"
)
