package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DealerClient struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DealerID primitive.ObjectID `json:"dealer_id" bson:"dealer_id"`
	PropertyID primitive.ObjectID `json:"property_id" bson:"property_id"`
	Name       string             `json:"name" bson:"name"`
	Phone      string             `json:"phone" bson:"phone"`
}