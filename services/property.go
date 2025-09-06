package services

import (
	"context"
	"fmt"
	"myapp/models"
	"myapp/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PropertyService struct {
	PropertyCollection *mongo.Collection
	CounterCollection  *mongo.Collection
}

func (s *PropertyService) CreateProperty(ctx context.Context, property models.Property) (primitive.ObjectID, error) {
	var resultID primitive.ObjectID

	err := utils.Retry(ctx, func() error {
		propertyNumber, err := s.getNextPropertyNumber(ctx)
		if err != nil {
			return err
		}

		property.PropertyNumber = propertyNumber

		result, err := s.PropertyCollection.InsertOne(ctx, property)
		if err != nil {
			return err
		}

		resultID = result.InsertedID.(primitive.ObjectID)
		return nil
	})

	if err != nil {
		return primitive.NilObjectID, err
	}

	return resultID, nil
}

// ← GENERATE next property number using counter collection
func (s *PropertyService) getNextPropertyNumber(ctx context.Context) (int64, error) {
	// ← USE findOneAndUpdate to atomically increment counter
	var result struct {
		Value int64 `bson:"value"`
	}

	err := s.CounterCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": "property_counter"},
		bson.M{"$inc": bson.M{"value": 1}},
		&options.FindOneAndUpdateOptions{
			Upsert:         &[]bool{true}[0],
			ReturnDocument: &[]options.ReturnDocument{options.After}[0],
		},
	).Decode(&result)

	if err != nil {
		return 0, err
	}

	return result.Value, nil
}

func (s *PropertyService) GetPropertyByNumber(ctx context.Context, propertyNumber int64) (*models.Property, error) {
	var property models.Property
	err := s.PropertyCollection.FindOne(ctx, bson.M{"property_number": propertyNumber, "is_deleted": false}).Decode(&property)
	if err != nil {
		return nil, err
	}
	return &property, nil
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

func (s *PropertyService) GetPropertiesByDealer(ctx context.Context, dealerID primitive.ObjectID, page, limit int) ([]models.Property, error) {
	// ← VALIDATE inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Default limit
	}

	filter := bson.M{
		"dealer_id":  dealerID,
		"is_deleted": bson.M{"$ne": true},
		"sold":       bson.M{"$ne": true},
	}

	// ← CALCULATE skip for pagination
	skip := (page - 1) * limit

	// ← PRODUCTION-READY options
	opts := options.Find().
		SetSort(bson.M{"_id": -1}). // Newest first
		SetSkip(int64(skip)).       // Pagination
		SetLimit(int64(limit)).     // Memory protection
		SetBatchSize(100)           // Network optimization

	cursor, err := s.PropertyCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query properties: %w", err)
	}
	defer cursor.Close(ctx)

	var properties []models.Property
	if err := cursor.All(ctx, &properties); err != nil {
		return nil, fmt.Errorf("failed to decode properties: %w", err)
	}

	return properties, nil
}

func (s *PropertyService) SearchProperties(ctx context.Context, filter bson.M, page, limit int, fields []string) ([]models.Property, error) {
	skip := (page - 1) * limit

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$skip", Value: int64(skip)}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := s.PropertyCollection.Aggregate(ctx, pipeline)
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
