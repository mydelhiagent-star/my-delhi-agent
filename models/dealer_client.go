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

type DealerClientPropertyInterest struct {
	ID         string        `json:"id"`
	PropertyID string        `json:"property_id"`
	Note       string        `json:"note"`
	Status     string        `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type DealerClientUpdate struct {
	Name       *string    `json:"name"`
	Phone      *string    `json:"phone"`
	Email      *string    `json:"email"`
	Note       *string    `json:"note"`
	UpdatedAt  *time.Time `json:"updated_at"`
	PropertyInterests *[]DealerClientPropertyInterestUpdate `json:"properties"`
}

type DealerClientPropertyInterestUpdate struct {
    ID         *string    `json:"id,omitempty"`
    PropertyID *string    `json:"property_id,omitempty"`
    Note       *string    `json:"note,omitempty"`
    Status     *string    `json:"status,omitempty"`
    UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}



type DealerClientQueryParams struct {
	ID              *string  `query:"id" mongo:"_id" convert:"objectid"`
	DealerID        *string  `query:"dealer_id" mongo:"dealer_id" convert:"objectid"`
	Name            *string  `query:"name"`
	Phone           *string  `query:"phone"`
	Note            *string  `query:"note"`

	PropertyInterestsID         *string `query:"properties_id" mongo:"properties._id" convert:"objectid"`
    PropertyInterestsPropertyID *string `query:"properties_property_id" mongo:"properties.property_id" convert:"objectid"`
    PropertyInterestsStatus     *string `query:"properties_status" mongo:"properties.status"`
    PropertyInterestsCreatedAt  *time.Time `query:"properties_created_at" mongo:"properties.created_at" convert:"date"`
    PropertyInterestsUpdatedAt  *time.Time `query:"properties_updated_at" mongo:"properties.updated_at" convert:"date"`
	
	CreatedAt       *time.Time `query:"created_at" mongo:"created_at" convert:"date"`
	UpdatedAt       *time.Time `query:"updated_at" mongo:"updated_at" convert:"date"`
	
	BaseQueryParams
}

// Dealer client status constants
const (
	DealerClientStatusMarked   = "marked"
	DealerClientStatusUnmarked = "unmarked"
)
