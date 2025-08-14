package models

import(
	 "go.mongodb.org/mongo-driver/bson/primitive"
)

type Property struct{
	ID              primitive.ObjectID      `json:"id,omitempty" bson:"_id,omitempty"`
	DealerID        primitive.ObjectID      `json:"dealer_id" bson:"dealer_id"`
	Title           string                  `json:"title" bson:"title"` 
	Description     string                  `json:"description" bson:"description"`  
	Address         string                  `json:"address" bson:"address"`
	Price           float64                 `json:"price" bson:"price"`
	Photos          []string                 `json:"photos,omitempty" bson:"photos,omitempty"`
    Videos          []string                 `json:"videos,omitempty" bson:"videos,omitempty"`
	OwnerName       string                  `json:"owner_name,omitempty" bson:"owner_name,omitempty"`
	OwnerPhone      string                  `json:"owner_phone,omitempty" bson:"owner_phone,omitempty"`
	NearestLandmark string                  `json:"nearest_landmark" bson:"nearest_landmark"`
}