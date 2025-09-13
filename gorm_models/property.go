package gorm_models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Property struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	PropertyNumber  int64          `gorm:"uniqueIndex;not null" json:"property_number"`
	DealerID        uint           `gorm:"not null;index" json:"dealer_id"`
	Dealer          Dealer         `gorm:"foreignKey:DealerID;references:ID" json:"dealer,omitempty"`
	Title           string         `gorm:"size:255;not null" json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	Address         string         `gorm:"type:text;not null" json:"address"`
	MinPrice        int64          `gorm:"not null" json:"min_price"`
	MaxPrice        int64          `gorm:"not null" json:"max_price"`
	Photos          pq.StringArray `gorm:"type:text[]" json:"photos"`
	Videos          pq.StringArray `gorm:"type:text[]" json:"videos"`
	OwnerName       string         `gorm:"size:255" json:"owner_name"`
	OwnerPhone      string         `gorm:"size:20" json:"owner_phone"`
	NearestLandmark string         `gorm:"size:255" json:"nearest_landmark"`
	IsDeleted       bool           `gorm:"default:false;index" json:"is_deleted"`
	Sold            bool           `gorm:"default:false;index" json:"sold"`
	SoldPrice       int64          `json:"sold_price"`
	SoldDate        *time.Time     `json:"sold_date"`
	Area            float64        `json:"area"`
	Bedrooms        int            `json:"bedrooms"`
	Bathrooms       int            `json:"bathrooms"`
	PropertyType    string         `gorm:"size:50;index" json:"property_type"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	LeadInterests []LeadPropertyInterest `gorm:"foreignKey:PropertyID;references:ID" json:"lead_interests,omitempty"`
	DealerClients []DealerClient         `gorm:"foreignKey:PropertyID;references:ID" json:"dealer_clients,omitempty"`
}

type PropertyUpdate struct {
	Title           *string    `json:"title,omitempty"`
	Address         *string    `json:"address,omitempty"`
	NearestLandmark *string    `json:"nearest_landmark,omitempty"`
	SoldBy          *string    `json:"sold_by,omitempty"`
	MinPrice        *int64     `json:"min_price,omitempty"`
	MaxPrice        *int64     `json:"max_price,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Photos          *[]string  `json:"photos,omitempty"`
	Videos          *[]string  `json:"videos,omitempty"`
	OwnerName       *string    `json:"owner_name,omitempty"`
	OwnerPhone      *string    `json:"owner_phone,omitempty"`
	Sold            *bool      `json:"sold,omitempty"`
	IsDeleted       *bool      `json:"is_deleted,omitempty"`
	SoldPrice       *int64     `json:"sold_price,omitempty"`
	SoldDate        *time.Time `json:"sold_date,omitempty"`
	Area            *float64   `json:"area,omitempty"`
	Bedrooms        *int       `json:"bedrooms,omitempty"`
	Bathrooms       *int       `json:"bathrooms,omitempty"`
	PropertyType    *string    `json:"property_type,omitempty"`
}
