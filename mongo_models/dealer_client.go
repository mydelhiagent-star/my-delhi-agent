package mongo_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DealerClient struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DealerID primitive.ObjectID `json:"dealer_id" bson:"dealer_id"`
	Name       string             `json:"name" bson:"name"`
	Phone      string             `json:"phone" bson:"phone"`
	Note       string             `json:"note" bson:"note"`
	PropertyInterests []DealerClientPropertyInterest `json:"properties" bson:"properties"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`

}

type DealerClientPropertyInterest struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DealerClientID primitive.ObjectID `json:"dealer_client_id" bson:"dealer_client_id"`
	PropertyID primitive.ObjectID `json:"property_id" bson:"property_id"`
	Status     string             `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}