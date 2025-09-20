package mongo_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Property struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	PropertyNumber  int64              `bson:"property_number"`
	DealerID        primitive.ObjectID `bson:"dealer_id"`
	Title           string             `bson:"title"`
	Description     string             `bson:"description"`
	Address         string             `bson:"address"`
	MinPrice        int64              `bson:"min_price"`
	MaxPrice        int64              `bson:"max_price"`
	Photos         []string           `bson:"photos,omitempty"`
	Videos          []string           `bson:"videos,omitempty"`
	OwnerName       string             `bson:"owner_name"`
	OwnerPhone      string             `bson:"owner_phone"`
	NearestLandmark string             `bson:"nearest_landmark"`
	IsDeleted       bool               `bson:"is_deleted,omitempty"`
	Sold            bool               `bson:"sold,omitempty"`
	SoldPrice       int64              `bson:"sold_price,omitempty"`
	SoldDate        time.Time          `bson:"sold_date,omitempty"`
	CreatedAt       time.Time          `bson:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at"`
	Area            float64            `bson:"area"`
	Bedrooms        int                `bson:"bedrooms"`
	Bathrooms       int                `bson:"bathrooms"`
	PropertyType    string             `bson:"property_type"`
}


type PropertyUpdate struct {
    Title           *string    `bson:"title"`
    Address         *string    `bson:"address"`
    NearestLandmark *string    `bson:"nearest_landmark"`
    SoldBy          *string    `bson:"sold_by"`
    MinPrice        *int64     `bson:"min_price"`
    MaxPrice        *int64     `bson:"max_price"`
    Description     *string    `bson:"description"`
    Photos          *[]string  `bson:"photos"`
    Videos          *[]string  `bson:"videos"`
    OwnerName       *string    `bson:"owner_name"`
    OwnerPhone      *string    `bson:"owner_phone"`
    Sold            *bool      `bson:"sold"`
    IsDeleted       *bool      `bson:"is_deleted"`
    SoldPrice       *int64     `bson:"sold_price"`
    SoldDate        *time.Time `bson:"sold_date"`
    Area            *float64   `bson:"area"`
    Bedrooms        *int       `bson:"bedrooms"`
    Bathrooms       *int       `bson:"bathrooms"`
    PropertyType    *string    `bson:"property_type"`
    UpdatedAt       *time.Time `bson:"updated_at"`
}


