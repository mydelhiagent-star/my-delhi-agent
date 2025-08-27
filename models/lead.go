package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lead struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Phone      string             `json:"phone" bson:"phone"`
	Properties []PropertyInterest `json:"properties,omitempty" bson:"properties,omitempty"`

	AadharNumber string `json:"aadhar_number,omitempty" bson:"aadhar_number,omitempty"`
	AadharPhoto  string `json:"aadhar_photo,omitempty" bson:"aadhar_photo,omitempty"`

	PopulatedProperties []Property `json:"populated_properties,omitempty" bson:"populated_properties,omitempty"`
}

type PropertyInterest struct {
	PropertyID primitive.ObjectID `json:"property_id" bson:"property_id"`
	DealerID   primitive.ObjectID `json:"dealer_id" bson:"dealer_id"` // ← ADD THIS
	Status     string             `json:"status" bson:"status"`
}

const (
	LeadStatusViewed       = "viewed"
	LeadStatusInterested   = "interested"
	LeadStatusBooked       = "booked"
	LeadStatusCancelled    = "cancelled"
	LeadStatusFailed       = "failed"
	LeadStatusUninterested = "uninterested"
)

// Add these structs to your existing models/lead.go
