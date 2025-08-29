package services

import (
	"context"
	"myapp/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Assume you have a GCPService with these methods:
// GenerateSignedURL(objectName string, expiry time.Duration) (string, error)
// PublicFileURL(objectName string) string

type PropertyService struct {
	PropertyCollection *mongo.Collection
}

func (s *PropertyService) CreateProperty(ctx context.Context, property models.Property) (primitive.ObjectID, error) {
	property.ID = primitive.NewObjectID()
	_, err := s.PropertyCollection.InsertOne(ctx, property)
	return property.ID, err
}

func (s *PropertyService) GetPropertyByID(id primitive.ObjectID) (*models.Property, error) {
	var property models.Property
	err := s.PropertyCollection.FindOne(context.Background(), bson.M{"_id": id, "is_deleted": false}).Decode(&property)
	if err != nil {
		return nil, err
	}
	return &property, nil
}

func (s *PropertyService) UpdateProperty(id primitive.ObjectID, updates models.PropertyUpdate) error {
	update := bson.M{"$set": updates}
	_, err := s.PropertyCollection.UpdateByID(context.Background(), id, update)
	return err
}

func (s *PropertyService) DeleteProperty(id primitive.ObjectID) error {
	_, err := s.PropertyCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"is_deleted": true}})
	return err
}

func (s *PropertyService) GetAllProperties(ctx context.Context) ([]models.Property, error) {
	cur, err := s.PropertyCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var properties []models.Property
	if err := cur.All(ctx, &properties); err != nil {
		return nil, err
	}
	return properties, nil
}

func (s *PropertyService) GetPropertiesByDealer(ctx context.Context, dealerID primitive.ObjectID) ([]models.Property, error) {
	filter := bson.M{"dealer_id": dealerID,"$or": []bson.M{
        {"is_deleted": false},           // Field exists and is false
        {"is_deleted": bson.M{"$exists": false}}, // Field doesn't exist
    },}

	cursor, err := s.PropertyCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return nil, err
	}

	return properties, nil
}
