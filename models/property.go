package models

import(
	 "go.mongodb.org/mongo-driver/bson/primitive"
)

type Property struct{
	ID             primitive.ObjectID      `json:"id,omitempty" bson:"_id,omitempty"`
	DealerID       primitive.ObjectID      `json:"dealer_id" bson:"dealer_id"`
	Title          string                  `json:"title" bson:"title"` 
	Description    string                  `json:"description" bson:"description"`  
	Area           string                  `json:"area" bson:"area"`
	Price          float64                 `json:"price" bson:"price"`
	Photos        []string                 `json:"photos" bson:"photos"`
    Videos        []string                 `json:"videos" bson:"videos"`
}