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
	PropertyNumber int64    `json:"property_number"`
	DealerID   string    `json:"dealer_id"`
	Status     string    `json:"status"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// models/lead.go
type LeadQueryParams struct {
	ID          *string    `query:"id" mongo:"_id" convert:"objectid"`
	Name        *string    `query:"name" mongo:"name"`
	Phone       *string    `query:"phone" mongo:"phone"`
	AadharNumber *string   `query:"aadhar_number" mongo:"aadhar_number"`
	DealerID    *string    `query:"dealer_id" mongo:"properties.dealer_id" convert:"objectid" array:"properties"`
	PropertyID  *string    `query:"property_id" mongo:"properties.property_id" convert:"objectid" array:"properties"`
	PropertyNumber *int64  `query:"property_number" mongo:"properties.property_number" convert:"int64" array:"properties"`
	Status      *string    `query:"status" mongo:"status"`
	CreatedAt   *time.Time `query:"created_at" mongo:"created_at" convert:"date"`
	UpdatedAt   *time.Time `query:"updated_at" mongo:"updated_at" convert:"date"`
	BaseQueryParams
}


const (
	LeadStatusView         = "view"
	LeadStatusInterested   = "interested"
	LeadStatusBooked       = "booked"
	LeadStatusCancelled    = "cancelled"
	LeadStatusFailed       = "failed"
	LeadStatusUninterested = "uninterested"
)
