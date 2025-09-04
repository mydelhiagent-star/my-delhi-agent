package services

import (
	"context"
	"myapp/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DealerClientService struct {
	DealerClientCollection *mongo.Collection
}

func (s *DealerClientService) CreateDealerClient(ctx context.Context, dealerClient models.DealerClient) (primitive.ObjectID, error) {
	_, err := s.DealerClientCollection.InsertOne(ctx, dealerClient)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return dealerClient.ID, nil
}
func (s *DealerClientService) GetDealerClientByPropertyID(ctx context.Context, dealerID primitive.ObjectID, propertyID primitive.ObjectID) ([]models.DealerClient, error) {
	cursor, err := s.DealerClientCollection.Find(ctx, bson.M{"dealer_id": dealerID, "property_id": propertyID})
	if err != nil {
		return nil, err
	}
	var dealerClients []models.DealerClient
	if err := cursor.All(ctx, &dealerClients); err != nil {
		return nil, err
	}
	return dealerClients, nil
}

func (s *DealerClientService) UpdateDealerClient(ctx context.Context, dealerClientID primitive.ObjectID, updateData interface{}) error {
	_, err := s.DealerClientCollection.UpdateByID(ctx, dealerClientID, bson.M{"$set": updateData})
	return err
}
