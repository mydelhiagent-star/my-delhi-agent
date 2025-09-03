package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Property struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PropertyNumber  int64              `json:"property_number" bson:"property_number"`
	DealerID        primitive.ObjectID `json:"dealer_id" bson:"dealer_id"`
	Title           string             `json:"title" bson:"title"`
	Description     string             `json:"description" bson:"description"`
	Address         string             `json:"address" bson:"address"`
	MinPrice        float64            `json:"min_price" bson:"min_price"`
	MaxPrice        float64            `json:"max_price" bson:"max_price"`
	Photos          []string           `json:"photos,omitempty" bson:"photos,omitempty"`
	Videos          []string           `json:"videos,omitempty" bson:"videos,omitempty"`
	OwnerName       string             `json:"owner_name,omitempty" bson:"owner_name,omitempty"`
	OwnerPhone      string             `json:"owner_phone,omitempty" bson:"owner_phone,omitempty"`
	NearestLandmark string             `json:"nearest_landmark" bson:"nearest_landmark"`
	IsDeleted       bool               `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	Sold            bool               `json:"sold,omitempty" bson:"sold,omitempty"`
	SoldPrice       float64            `json:"sold_price,omitempty" bson:"sold_price,omitempty"`
}

type PropertyUpdate struct {
	Title           *string   `json:"title,omitempty" bson:"title,omitempty"`
	Address         *string   `json:"address,omitempty" bson:"address,omitempty"`
	NearestLandmark *string   `json:"nearest_landmark,omitempty" bson:"nearest_landmark,omitempty"`
	SoldBy          *string   `json:"sold_by,omitempty" bson:"sold_by,omitempty"`
	MinPrice        *float64  `json:"min_price,omitempty" bson:"min_price,omitempty"`
	MaxPrice        *float64  `json:"max_price,omitempty" bson:"max_price,omitempty"`
	Description     *string   `json:"description,omitempty" bson:"description,omitempty"`
	Photos          *[]string `json:"photos,omitempty" bson:"photos,omitempty"`
	Videos          *[]string `json:"videos,omitempty" bson:"videos,omitempty"`
	OwnerName       *string   `json:"owner_name,omitempty" bson:"owner_name,omitempty"`
	OwnerPhone      *string   `json:"owner_phone,omitempty" bson:"owner_phone,omitempty"`
	Sold            *bool     `json:"sold,omitempty" bson:"sold,omitempty"`
	IsDeleted       *bool     `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	SoldPrice       *float64  `json:"sold_price,omitempty" bson:"sold_price,omitempty"`
}
