package mongo_repository

import (
	"context"
	"errors"
	"fmt"
	"myapp/mongo_models"
	"myapp/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoLeadRepository struct {
	leadCollection     *mongo.Collection
	propertyCollection *mongo.Collection
}

func NewMongoLeadRepository(leadCollection, propertyCollection *mongo.Collection) repositories.LeadRepository {
	return &MongoLeadRepository{
		leadCollection:     leadCollection,
		propertyCollection: propertyCollection,
	}
}

func (r *MongoLeadRepository) Create(ctx context.Context, lead models.Lead) (primitive.ObjectID, error) {
	res, err := r.leadCollection.InsertOne(ctx, lead)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("failed to convert inserted ID")
	}
	return id, nil
}

func (r *MongoLeadRepository) GetByID(ctx context.Context, id primitive.ObjectID) (models.Lead, error) {
	var lead models.Lead
	err := r.leadCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&lead)
	return lead, err
}

func (r *MongoLeadRepository) GetAll(ctx context.Context) ([]models.Lead, error) {
	cursor, err := r.leadCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var leads []models.Lead
	if err = cursor.All(ctx, &leads); err != nil {
		return nil, err
	}
	return leads, nil
}

func (r *MongoLeadRepository) GetByDealerID(ctx context.Context, dealerID primitive.ObjectID) ([]models.Lead, error) {
	cursor, err := r.leadCollection.Find(ctx, bson.M{"dealer_id": dealerID})
	if err != nil {
		return nil, err
	}

	var leads []models.Lead
	if err = cursor.All(ctx, &leads); err != nil {
		return nil, err
	}
	return leads, nil
}

func (r *MongoLeadRepository) Search(ctx context.Context, filter bson.M, page, limit int, fields []string) ([]models.Lead, error) {
	skip := (page - 1) * limit

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$skip", Value: int64(skip)}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := r.leadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leads []models.Lead
	if err := cursor.All(ctx, &leads); err != nil {
		return nil, err
	}
	return leads, nil
}

func (r *MongoLeadRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	update := bson.M{"$set": updates}
	_, err := r.leadCollection.UpdateByID(ctx, id, update)
	return err
}

func (r *MongoLeadRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.leadCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoLeadRepository) AddPropertyInterest(ctx context.Context, leadID primitive.ObjectID, propertyInterest models.PropertyInterest) error {
	// Check if property already exists for this lead
	filter := bson.M{
		"_id":                    leadID,
		"properties.property_id": propertyInterest.PropertyID,
	}

	var existingLead models.Lead
	err := r.leadCollection.FindOne(ctx, filter).Decode(&existingLead)
	if err == nil {
		return errors.New("property already added to this lead")
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("database error checking property: %w", err)
	}

	// Add property interest
	update := bson.M{
		"$push": bson.M{
			"properties": propertyInterest,
		},
	}

	_, err = r.leadCollection.UpdateOne(ctx, bson.M{"_id": leadID}, update)
	return err
}

func (r *MongoLeadRepository) UpdatePropertyInterest(ctx context.Context, leadID, propertyID primitive.ObjectID, status, note string) error {
	filter := bson.M{
		"_id":                    leadID,
		"properties.property_id": propertyID,
	}

	update := bson.M{
		"$set": bson.M{
			"properties.$.status": status,
			"properties.$.note":   note,
		},
	}

	result, err := r.leadCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *MongoLeadRepository) GetLeadPropertyDetails(ctx context.Context, leadID primitive.ObjectID) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": leadID}}},
		{{Key: "$unwind", Value: "$properties"}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "property",
			"localField":   "properties.property_id",
			"foreignField": "_id",
			"as":           "property_details",
		}}},
		{{Key: "$unwind", Value: "$property_details"}},
		{{Key: "$project", Value: bson.M{
			"property_id":   "$properties.property_id",
			"dealer_id":     "$properties.dealer_id",
			"status":        "$properties.status",
			"note":          "$properties.note",
			"created_at":    "$properties.created_at",
			"title":         "$property_details.title",
			"address":       "$property_details.address",
			"min_price":     "$property_details.min_price",
			"max_price":     "$property_details.max_price",
			"photos":        "$property_details.photos",
			"property_type": "$property_details.property_type",
			"bedrooms":      "$property_details.bedrooms",
			"bathrooms":     "$property_details.bathrooms",
		}}},
	}

	cursor, err := r.leadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *MongoLeadRepository) GetDealerLeads(ctx context.Context, dealerID primitive.ObjectID) ([]models.Lead, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"properties.dealer_id": dealerID,
		}}},
		{{Key: "$addFields", Value: bson.M{
			"properties": bson.M{
				"$filter": bson.M{
					"input": "$properties",
					"cond":  bson.M{"$eq": []interface{}{"$$this.dealer_id", dealerID}},
				},
			},
		}}},
	}

	cursor, err := r.leadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leads []models.Lead
	if err := cursor.All(ctx, &leads); err != nil {
		return nil, err
	}

	return leads, nil
}

func (r *MongoLeadRepository) GetPropertyDetails(ctx context.Context, soldStr, deletedStr string) ([]bson.M, error) {
	filter := bson.M{}

	if soldStr != "" {
		if soldStr == "true" {
			filter["sold"] = true
		} else if soldStr == "false" {
			filter["sold"] = bson.M{"$ne": true}
		}
	}

	if deletedStr != "" {
		if deletedStr == "true" {
			filter["is_deleted"] = true
		} else if deletedStr == "false" {
			filter["is_deleted"] = bson.M{"$ne": true}
		}
	}

	cursor, err := r.propertyCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *MongoLeadRepository) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	count, err := r.leadCollection.CountDocuments(ctx, bson.M{"phone": phone})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
