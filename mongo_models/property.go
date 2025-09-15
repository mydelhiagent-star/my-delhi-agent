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
	UpdatedAt       time.Time          `bson:"updated_at,omitempty"`
	Area            float64            `bson:"area"`
	Bedrooms        int                `bson:"bedrooms"`
	Bathrooms       int                `bson:"bathrooms"`
	PropertyType    string             `bson:"property_type"`
}

type PropertyUpdate struct {
	Title           *string    `bson:"title,omitempty"`
	Address         *string    `bson:"address,omitempty"`
	NearestLandmark *string    `bson:"nearest_landmark,omitempty"`
	SoldBy          *string    `bson:"sold_by,omitempty"`
	MinPrice        *int64     `bson:"min_price,omitempty"`
	MaxPrice        *int64     `bson:"max_price,omitempty"`
	Description     *string    `bson:"description,omitempty"`
	Photos          *[]string  `bson:"photos,omitempty"`
	Videos          *[]string  `bson:"videos,omitempty"`
	OwnerName       *string    `bson:"owner_name,omitempty"`
	OwnerPhone      *string    `bson:"owner_phone,omitempty"`
	Sold            *bool      `bson:"sold,omitempty"`
	IsDeleted       *bool      `bson:"is_deleted,omitempty"`
	SoldPrice       *int64     `bson:"sold_price,omitempty"`
	SoldDate        *time.Time `bson:"sold_date,omitempty"`
	Area            *float64   `bson:"area,omitempty"`
	Bedrooms        *int       `bson:"bedrooms,omitempty"`
	Bathrooms       *int       `bson:"bathrooms,omitempty"`
	PropertyType    *string    `bson:"property_type,omitempty"`
	UpdatedAt       *time.Time `bson:"updated_at,omitempty"`
}
