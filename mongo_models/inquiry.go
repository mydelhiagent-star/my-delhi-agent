package mongo_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inquiry struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	DealerID    *primitive.ObjectID `bson:"dealer_id,omitempty"`
	Source      string              `bson:"source"`
	Name        string              `bson:"name"`
	Phone       string              `bson:"phone"`
	Requirement string              `bson:"requirement"`
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
}

type InquiryUpdate struct {
	DealerID    *primitive.ObjectID `bson:"dealer_id"`
	Source      *string             `bson:"source"`
	Name        *string             `bson:"name"`
	Phone       *string             `bson:"phone"`
	Requirement *string             `bson:"requirement"`
	UpdatedAt   *time.Time          `bson:"updated_at"`
}
