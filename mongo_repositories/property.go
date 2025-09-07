package mongo_repository

import (
	"context"
	"fmt"
	"myapp/models"
	"myapp/repositories"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPropertyRepository struct {
	propertyCollection *mongo.Collection
	counterCollection  *mongo.Collection
	redisClient        *redis.Client
}

func NewMongoPropertyRepository(propertyCollection, counterCollection *mongo.Collection, redisClient *redis.Client) repositories.PropertyRepository {
	return &MongoPropertyRepository{
		propertyCollection: propertyCollection,
		counterCollection:  counterCollection,
		redisClient:        redisClient,
	}
}

func (r *MongoPropertyRepository) Create(ctx context.Context, property models.Property) (primitive.ObjectID, error) {
	result, err := r.propertyCollection.InsertOne(ctx, property)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *MongoPropertyRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Property, error) {
	var property models.Property
	err := r.propertyCollection.FindOne(ctx, bson.M{"_id": id, "is_deleted": false}).Decode(&property)
	if err != nil {
		return nil, err
	}
	return &property, nil
}

func (r *MongoPropertyRepository) GetByNumber(ctx context.Context, propertyNumber int64) (*models.Property, error) {
	var property models.Property
	err := r.propertyCollection.FindOne(ctx, bson.M{"property_number": propertyNumber, "is_deleted": false}).Decode(&property)
	if err != nil {
		return nil, err
	}
	return &property, nil
}

func (r *MongoPropertyRepository) GetByDealer(ctx context.Context, dealerID primitive.ObjectID, page, limit int) ([]models.Property, error) {
	// Validate inputs
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

	// Calculate skip for pagination
	skip := (page - 1) * limit

	// Production-ready options
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}). // Newest first
		SetSkip(int64(skip)).              // Pagination
		SetLimit(int64(limit)).            // Memory protection
		SetBatchSize(100)                  // Network optimization

	cursor, err := r.propertyCollection.Find(ctx, filter, opts)
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

func (r *MongoPropertyRepository) GetAll(ctx context.Context) ([]models.Property, error) {
	cursor, err := r.propertyCollection.Find(ctx, bson.M{})
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

func (r *MongoPropertyRepository) Search(ctx context.Context, filter bson.M, page, limit int, fields []string) ([]models.Property, error) {
	skip := (page - 1) * limit

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$skip", Value: int64(skip)}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := r.propertyCollection.Aggregate(ctx, pipeline)
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

func (r *MongoPropertyRepository) Update(ctx context.Context, id primitive.ObjectID, updates models.PropertyUpdate) error {
	update := bson.M{"$set": updates}
	_, err := r.propertyCollection.UpdateByID(ctx, id, update)
	return err
}

func (r *MongoPropertyRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	_, err := r.propertyCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"is_deleted": true,
		"updated_at": now,
	}})
	return err
}

func (r *MongoPropertyRepository) GetNextPropertyNumber(ctx context.Context) (int64, error) {
	// Use findOneAndUpdate to atomically increment counter
	var result struct {
		Value int64 `bson:"value"`
	}

	err := r.counterCollection.FindOneAndUpdate(
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
