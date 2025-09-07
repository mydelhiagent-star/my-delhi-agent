package services

import (
	"context"
	"errors"
	"myapp/models"
	"myapp/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DealerClientService struct {
	Repo repositories.DealerClientRepository
}

func (s *DealerClientService) CheckPhoneExistsForDealer(ctx context.Context, dealerID primitive.ObjectID, propertyID primitive.ObjectID, phone string) (bool, error) {
	// This would need a custom method in the repository
	// For now, return false
	return false, nil
}

func (s *DealerClientService) CreateDealerClient(ctx context.Context, dealerClient models.DealerClient) (primitive.ObjectID, error) {
	// Check if phone number already exists for this dealer
	exists, err := s.CheckPhoneExistsForDealer(ctx, dealerClient.DealerID, dealerClient.PropertyID, dealerClient.Phone)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if exists {
		return primitive.NilObjectID, errors.New("phone number already exists")
	}

	// Set default status
	dealerClient.Status = "unmarked"

	return s.Repo.Create(ctx, dealerClient)
}

func (s *DealerClientService) GetDealerClientByPropertyID(ctx context.Context, dealerID primitive.ObjectID, propertyID primitive.ObjectID) ([]models.DealerClient, error) {
	return s.Repo.GetByPropertyID(ctx, propertyID)
}

func (s *DealerClientService) GetDealerClientsByDealerID(ctx context.Context, dealerID primitive.ObjectID) ([]models.DealerClient, error) {
	return s.Repo.GetByDealerID(ctx, dealerID)
}

func (s *DealerClientService) GetAllDealerClients(ctx context.Context) ([]models.DealerClient, error) {
	return s.Repo.GetAll(ctx)
}

func (s *DealerClientService) UpdateDealerClient(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return s.Repo.Update(ctx, id, updates)
}

func (s *DealerClientService) DeleteDealerClient(ctx context.Context, id primitive.ObjectID) error {
	return s.Repo.Delete(ctx, id)
}

func (s *DealerClientService) GetDealerClientByID(ctx context.Context, id primitive.ObjectID) (*models.DealerClient, error) {
	return s.Repo.GetByID(ctx, id)
}
