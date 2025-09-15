package models

import "time"

type DealerClient struct {
	ID         string    `json:"id"`
	DealerID   string    `json:"dealer_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	PropertyInterests []DealerClientPropertyInterest `json:"properties"`
}

type DealerClientUpdate struct {
	Name       *string    `json:"name"`
	Phone      *string    `json:"phone"`
	Email      *string    `json:"email"`
	Note       *string    `json:"note"`
	UpdatedAt  *time.Time `json:"updated_at"`
	PropertyInterests *[]DealerClientPropertyInterest `json:"properties"`
}


type DealerClientPropertyInterest struct {
	ID         string        `json:"id"`
	PropertyID string        `json:"property_id"`
	Status     string        `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// Dealer client status constants
const (
	DealerClientStatusMarked   = "marked"
	DealerClientStatusUnmarked = "unmarked"
)
