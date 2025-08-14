package models

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lead struct{
	ID                primitive.ObjectID      `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string                  `json:"name" bson:"name"`
	Phone			  string                  `json:"phone" bson:"phone"`
	Requirement       string                  `json:"requirement" bson:"requirement"`
	Status            string                  `json:"status" bson:"status"`
	AadharNumber      string                  `json:"aadhar_number,omitempty" bson:"aadhar_number,omitempty"`
	AadharPhoto       string                  `json:"aadhar_photo,omitempty" bson:"aadhar_photo,omitempty"`
}

const (
	LeadStatusNew              = "new"
	LeadStatusInProgress       = "in_progress"
	LeadStatusSuccess          = "success"
)


