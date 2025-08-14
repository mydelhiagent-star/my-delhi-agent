package models

import(
	 "go.mongodb.org/mongo-driver/bson/primitive"
)
type Dealer struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Phone    string `json:"phone" bson:"phone"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	OfficeAddress string `json:"office_address" bson:"office_address"`
	ShopName string `json:"shop_name" bson:"shop_name"`
	Location string `json:"location" bson:"location"`
	SubLocation string `json:"sub_location" bson:"sub_location"`
}
