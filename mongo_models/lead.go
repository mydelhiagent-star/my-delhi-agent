package mongo_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lead struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Phone       string             `json:"phone" bson:"phone"`
	Requirement string             `json:"requirement" bson:"requirement"`
	Properties  []PropertyInterest `json:"properties,omitempty" bson:"properties,omitempty"`

	AadharNumber string `json:"aadhar_number,omitempty" bson:"aadhar_number,omitempty"`
	AadharPhoto  string `json:"aadhar_photo,omitempty" bson:"aadhar_photo,omitempty"`

	PopulatedProperties []Property `json:"populated_properties,omitempty" bson:"populated_properties,omitempty"`
}

type PropertyInterest struct {
	PropertyID primitive.ObjectID `json:"property_id" bson:"property_id"`
	PropertyNumber int64    `json:"property_number" bson:"property_number"`
	DealerID   primitive.ObjectID `json:"dealer_id" bson:"dealer_id"` // ← ADD THIS
	Status     string             `json:"status" bson:"status"`
	Note       string             `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

const (
	LeadStatusView         = "view"
	LeadStatusInterested   = "interested"
	LeadStatusBooked       = "booked"
	LeadStatusCancelled    = "cancelled"
	LeadStatusFailed       = "failed"
	LeadStatusUninterested = "uninterested"
)

// Add these structs to your existing models/lead.go
