package mongo_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inquiry struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Phone       string             `bson:"phone"`
	Requirement string             `bson:"requirement"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type InquiryUpdate struct {
	Name        *string    `bson:"name"`
	Phone       *string    `bson:"phone"`
	Requirement *string    `bson:"requirement"`
	UpdatedAt   *time.Time `bson:"updated_at"`
}
