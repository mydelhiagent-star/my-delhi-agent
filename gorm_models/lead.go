package gorm_models

import (
	"time"

	"gorm.io/gorm"
)

type Lead struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Phone        string         `gorm:"size:20;not null;index" json:"phone"`
	Requirement  string         `gorm:"type:text" json:"requirement"`
	AadharNumber string         `gorm:"size:20" json:"aadhar_number"`
	AadharPhoto  string         `gorm:"size:255" json:"aadhar_photo"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	PropertyInterests []LeadPropertyInterest `gorm:"foreignKey:LeadID;references:ID" json:"property_interests,omitempty"`
}

type LeadPropertyInterest struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	LeadID     uint      `gorm:"not null;index" json:"lead_id"`
	PropertyID uint      `gorm:"not null;index" json:"property_id"`
	DealerID   uint      `gorm:"not null;index" json:"dealer_id"`
	Status     string    `gorm:"size:50;default:'view'" json:"status"`
	Note       string    `gorm:"type:text" json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relationships
	Lead     Lead     `gorm:"foreignKey:LeadID;references:ID" json:"lead,omitempty"`
	Property Property `gorm:"foreignKey:PropertyID;references:ID" json:"property,omitempty"`
	Dealer   Dealer   `gorm:"foreignKey:DealerID;references:ID" json:"dealer,omitempty"`
}

// Lead status constants
const (
	LeadStatusView         = "view"
	LeadStatusInterested   = "interested"
	LeadStatusBooked       = "booked"
	LeadStatusCancelled    = "cancelled"
	LeadStatusFailed       = "failed"
	LeadStatusUninterested = "uninterested"
)
