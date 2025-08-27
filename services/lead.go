package services

import (
	"context"
	"errors"

	"myapp/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LeadService struct {
	LeadCollection *mongo.Collection
}

func (s *LeadService) CreateLead(ctx context.Context, lead models.Lead) (primitive.ObjectID, error) {

	existingLead := models.Lead{}
	err := s.LeadCollection.FindOne(ctx, bson.M{"phone": lead.Phone}).Decode(&existingLead)
	if err == nil {
		return primitive.NilObjectID, errors.New("lead with phone already exists")
	} else if err != mongo.ErrNoDocuments {
		return primitive.NilObjectID, errors.New("database error checking phone")
	}

	res, err := s.LeadCollection.InsertOne(ctx, lead)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, err
	}
	return id, nil
}

func (s *LeadService) GetLeadByID(ctx context.Context, id primitive.ObjectID) (models.Lead, error) {

	var lead models.Lead
	err := s.LeadCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&lead)
	return lead, err
}

func (s *LeadService) GetAllLeads(ctx context.Context) ([]models.Lead, error) {
	cursor, err := s.LeadCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var leads []models.Lead
	if err = cursor.All(ctx, &leads); err != nil {
		return nil, err
	}
	return leads, nil
}

func (s *LeadService) GetAllLeadsByDealerID(ctx context.Context, dealerID primitive.ObjectID) ([]models.Lead, error) {
	cursor, err := s.LeadCollection.Find(ctx, bson.M{"dealer_id": dealerID})
	if err != nil {
		return nil, err
	}

	var leads []models.Lead
	if err = cursor.All(ctx, &leads); err != nil {
		return nil, err
	}
	return leads, nil

}

func (s *LeadService) UpdateLead(ctx context.Context, id primitive.ObjectID, updateData map[string]interface{}) error {

	update := bson.M{
		"$set": updateData,
	}
	_, err := s.LeadCollection.UpdateByID(ctx, id, update)
	return err
}

func (s *LeadService) AddPropertyInterest(ctx context.Context, leadID primitive.ObjectID, propertyInterest models.PropertyInterest) error {
	// Set timestamps and status

	propertyInterest.Status = models.LeadStatusViewed

	_, err := s.LeadCollection.UpdateOne(ctx,
		bson.M{"_id": leadID},
		bson.M{"$addToSet": bson.M{"properties": propertyInterest}},
	)

	return err
}

func (s *LeadService) SearchLeads(ctx context.Context, filter bson.M, page, limit int) ([]models.Lead, error) {
	// ← CALCULATE skip value for pagination
	skip := (page - 1) * limit

	// ← GET total count first

	// ← SET pagination options
	options := options.Find()
	options.SetSort(bson.D{{Key: "_id", Value: -1}}) // Newest first
	options.SetLimit(int64(limit))
	options.SetSkip(int64(skip))

	cursor, err := s.LeadCollection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leads []models.Lead
	if err = cursor.All(ctx, &leads); err != nil {
		return nil, err
	}

	return leads, nil
}

func (s *LeadService) GetLeadPropertyDetails(ctx context.Context, leadID primitive.ObjectID) (*[]models.Property, error) {
	// ← BUILD the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match the specific lead
		{{Key: "$match", Value: bson.M{"_id": leadID}}},

		// Stage 2: Lookup properties
		{{Key: "$lookup", Value: bson.M{
			"from":         "property",
			"localField":   "properties.property_id",
			"foreignField": "_id",
			"as":           "populated_properties",
		}}},
		
	}

	// ← EXECUTE the aggregation pipeline
	cursor, err := s.LeadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// ← DECODE the result
	var result models.Lead
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	} else {
		return nil, mongo.ErrNoDocuments
	}

	return &result.PopulatedProperties, nil
}
