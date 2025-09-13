package models

import "time"

type Lead struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Phone        string                 `json:"phone"`
	Requirement  string                 `json:"requirement"`
	AadharNumber string                 `json:"aadhar_number"`
	AadharPhoto  string                 `json:"aadhar_photo"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Properties   []PropertyInterest `json:"properties,omitempty"`
}


type PropertyInterest struct {
	ID         string    `json:"id"`
	LeadID     string    `json:"lead_id"`
	PropertyID string    `json:"property_id"`
	DealerID   string    `json:"dealer_id"`
	Status     string    `json:"status"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}


const (
	LeadStatusView         = "view"
	LeadStatusInterested   = "interested"
	LeadStatusBooked       = "booked"
	LeadStatusCancelled    = "cancelled"
	LeadStatusFailed       = "failed"
	LeadStatusUninterested = "uninterested"
)
