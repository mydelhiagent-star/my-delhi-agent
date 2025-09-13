package models

import "time"

type DealerClient struct {
	ID         string    `json:"id"`
	DealerID   string    `json:"dealer_id"`
	PropertyID string    `json:"property_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Status     string    `json:"status"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Dealer client status constants
const (
	DealerClientStatusActive   = "active"
	DealerClientStatusInactive = "inactive"
	DealerClientStatusBlocked  = "blocked"
)
