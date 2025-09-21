package mongo_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	URL  string `bson:"url"`
	Type string `bson:"type"`
	Name string `bson:"name"`
	Size int64  `bson:"size"`
}

type DealerClient struct {
	ID                primitive.ObjectID             `bson:"_id,omitempty"`
	DealerID          primitive.ObjectID             `bson:"dealer_id"`
	Name              string                         `bson:"name"`
	Phone             string                         `bson:"phone"`
	Note              string                         `bson:"note"`
	Docs              []Document                     `bson:"docs,omitempty"`
	PropertyInterests []DealerClientPropertyInterest `bson:"properties,omitempty"`
	CreatedAt         time.Time                      `bson:"created_at"`
	UpdatedAt         time.Time                      `bson:"updated_at"`
}

type DealerClientPropertyInterest struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	PropertyID primitive.ObjectID `bson:"property_id"`
	Note       string             `bson:"note,omitempty"`
	Status     string             `bson:"status"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

type DealerClientUpdate struct {
	Name  *string     `bson:"name"`
	Phone *string     `bson:"phone"`
	Note  *string     `bson:"note"`
	Docs  *[]Document `bson:"docs"`
}
