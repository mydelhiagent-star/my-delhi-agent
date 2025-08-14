package services

import (
	"context"
	

	"myapp/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeadService struct {
	LeadCollection *mongo.Collection
}

func (s *LeadService) CreateLead(ctx context.Context,lead models.Lead) (primitive.ObjectID, error) {
	

	lead.Status = models.LeadStatusNew

	

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

func (s *LeadService) GetLeadByID(ctx context.Context,id primitive.ObjectID) (models.Lead, error) {
	

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

func (s *LeadService) GetAllLeadsByDealerID(ctx context.Context,dealerID primitive.ObjectID) ([]models.Lead, error) {
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

func (s *LeadService) UpdateLead(ctx context.Context,id primitive.ObjectID, updateData map[string]interface{}) error {
	
	update := bson.M{
		"$set": updateData,
	}
	_, err := s.LeadCollection.UpdateByID(ctx, id, update)
	return err
}
