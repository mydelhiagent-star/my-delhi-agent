package models

import "time"

// PropertyData represents a database-agnostic property model
type Property struct {
	ID              string    `json:"id"`
	PropertyNumber  int64     `json:"property_number"`
	DealerID        string    `json:"dealer_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Address         string    `json:"address"`
	MinPrice        int64     `json:"min_price"`
	MaxPrice        int64     `json:"max_price"`
	Photos          []string  `json:"photos"`
	Videos          []string  `json:"videos"`
	OwnerName       string    `json:"owner_name"`
	OwnerPhone      string    `json:"owner_phone"`
	NearestLandmark string    `json:"nearest_landmark"`
	IsDeleted       bool      `json:"is_deleted"`
	Sold            bool      `json:"sold"`
	SoldPrice       int64     `json:"sold_price"`
	SoldDate        time.Time `json:"sold_date"`
	Area            float64   `json:"area"`
	Bedrooms        int       `json:"bedrooms"`
	Bathrooms       int       `json:"bathrooms"`
	PropertyType    string    `json:"property_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
