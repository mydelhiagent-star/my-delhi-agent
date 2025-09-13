package mongo_repositories

import (
	"context"
	"fmt"
	"myapp/converters"
	"myapp/models"
	mongoModels "myapp/mongo_models"
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

func (r *MongoPropertyRepository) Create(ctx context.Context, property models.Property) (string, error) {
	mongoProperty, err := converters.ToMongoProperty(property)
	if err != nil {
		return "", err
	}

	result, err := r.propertyCollection.InsertOne(ctx, mongoProperty)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoPropertyRepository) GetByID(ctx context.Context, id string) (models.Property, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Property{}, err
	}

	var mongoProperty mongoModels.Property
	err = r.propertyCollection.FindOne(ctx, bson.M{"_id": objectID, "is_deleted": bson.M{"$ne": true}}).Decode(&mongoProperty)
	if err != nil {
		return models.Property{}, err
	}

	return converters.ToDomainProperty(mongoProperty), nil
}

func (r *MongoPropertyRepository) GetByDealer(ctx context.Context, dealerID string, page, limit int) ([]models.Property, error) {
	// Validate inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	dealerObjectID, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"dealer_id":  dealerObjectID,
		"is_deleted": bson.M{"$ne": true},
		"sold":       bson.M{"$ne": true},
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetBatchSize(100)

	cursor, err := r.propertyCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query properties: %w", err)
	}
	defer cursor.Close(ctx)

	var mongoProperties []mongoModels.Property
	if err := cursor.All(ctx, &mongoProperties); err != nil {
		return nil, fmt.Errorf("failed to decode properties: %w", err)
	}

	return converters.ToDomainPropertySlice(mongoProperties), nil
}

func (r *MongoPropertyRepository) GetAll(ctx context.Context) ([]models.Property, error) {
	cursor, err := r.propertyCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoProperties []mongoModels.Property
	if err := cursor.All(ctx, &mongoProperties); err != nil {
		return nil, err
	}

	return converters.ToDomainPropertySlice(mongoProperties), nil
}

func (r *MongoPropertyRepository) Update(ctx context.Context, id string, updates models.PropertyUpdate) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Convert PropertyUpdate to bson.M
	updateDoc := bson.M{}
	if updates.Title != nil {
		updateDoc["title"] = *updates.Title
	}
	if updates.Address != nil {
		updateDoc["address"] = *updates.Address
	}
	if updates.MinPrice != nil {
		updateDoc["min_price"] = *updates.MinPrice
	}
	if updates.MaxPrice != nil {
		updateDoc["max_price"] = *updates.MaxPrice
	}
	if updates.Sold != nil {
		updateDoc["sold"] = *updates.Sold
	}
	// Add other fields as needed...

	update := bson.M{"$set": updateDoc}
	_, err = r.propertyCollection.UpdateByID(ctx, objectID, update)
	return err
}

func (r *MongoPropertyRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = r.propertyCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{
		"is_deleted": true,
		"updated_at": now,
	}})
	return err
}

func (r *MongoPropertyRepository) GetNextPropertyNumber(ctx context.Context) (int64, error) {
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

func (r *MongoPropertyRepository) GetByNumber(ctx context.Context, propertyNumber int64) (models.Property, error) {
	var mongoProperty mongoModels.Property
	err := r.propertyCollection.FindOne(ctx, bson.M{"property_number": propertyNumber, "is_deleted": false}).Decode(&mongoProperty)
	if err != nil {
		return models.Property{}, err
	}

	return converters.ToDomainProperty(mongoProperty), nil
}

func (r *MongoPropertyRepository) Search(ctx context.Context, filter map[string]interface{}, page, limit int, fields []string) ([]models.Property, error) {
	skip := (page - 1) * limit

	// Convert map[string]interface{} to bson.M
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bsonFilter}},
		{{Key: "$skip", Value: int64(skip)}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := r.propertyCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoProperties []mongoModels.Property
	if err := cursor.All(ctx, &mongoProperties); err != nil {
		return nil, err
	}

	return converters.ToDomainPropertySlice(mongoProperties), nil
}
