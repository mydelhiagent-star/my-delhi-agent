package models

import "time"


type Dealer struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	OfficeAddress string    `json:"office_address"`
	ShopName      string    `json:"shop_name"`
	Location      string    `json:"location"`
	SubLocation   string    `json:"sub_location"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type LocationWithSubLocations struct {
	Location    string   `json:"location"`
	SubLocation []string `json:"sub_location"`
}
