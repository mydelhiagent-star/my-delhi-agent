package models

import "time"

type Inquiry struct {
	ID          string    `json:"id"`
	DealerID    *string   `json:"dealer_id,omitempty"`
	Source      string    `json:"source"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Requirement string    `json:"requirement"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InquiryUpdate struct {
	DealerID    *string    `json:"dealer_id,omitempty"`
	Source      *string    `json:"source,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Requirement *string    `json:"requirement,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type InquiryQueryParams struct {
	ID          *string `query:"id" mongo:"_id" convert:"objectid"`
	DealerID    *string `query:"dealer_id" mongo:"dealer_id" convert:"objectid"`
	Source      *string `query:"source"`
	Name        *string `query:"name"`
	Phone       *string `query:"phone"`
	Requirement *string `query:"requirement"`
	BaseQueryParams
}

func (i *InquiryQueryParams) SetDefaults() {
	// Call parent defaults
	i.BaseQueryParams.SetDefaults()
}
