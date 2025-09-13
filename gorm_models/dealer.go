package gorm_models

import (
	"time"

	"gorm.io/gorm"
)

type Dealer struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	Phone         string         `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	Password      string         `gorm:"size:255;not null" json:"password"`
	Email         string         `gorm:"size:255;uniqueIndex" json:"email"`
	OfficeAddress string         `gorm:"type:text" json:"office_address"`
	ShopName      string         `gorm:"size:255" json:"shop_name"`
	Location      string         `gorm:"size:100;index" json:"location"`
	SubLocation   string         `gorm:"size:100;index" json:"sub_location"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Properties    []Property             `gorm:"foreignKey:DealerID;references:ID" json:"properties,omitempty"`
	LeadInterests []LeadPropertyInterest `gorm:"foreignKey:DealerID;references:ID" json:"lead_interests,omitempty"`
	DealerClients []DealerClient         `gorm:"foreignKey:DealerID;references:ID" json:"dealer_clients,omitempty"`
}

type LocationWithSubLocations struct {
	Location    string   `json:"location"`
	SubLocation []string `json:"sub_location"`
}
