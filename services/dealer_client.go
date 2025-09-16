package services

import (
	"context"
	"errors"
	"myapp/models"
	"myapp/repositories"
)

type DealerClientService struct {
	Repo repositories.DealerClientRepository
}

func (s *DealerClientService) CheckPhoneExistsForDealer(ctx context.Context, dealerID string, phone string) (bool, error) {
	return s.Repo.CheckPhoneExistsForDealer(ctx, dealerID, phone)
}

func (s *DealerClientService) CreateDealerClient(ctx context.Context, dealerClient models.DealerClient) (string, error) {
	
	exists, err := s.CheckPhoneExistsForDealer(ctx, dealerClient.DealerID, dealerClient.Phone)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("phone number already exists")
	}

	

	return s.Repo.Create(ctx, dealerClient)
}

func (s *DealerClientService) GetDealerClients(ctx context.Context, params models.DealerClientQueryParams) ([]models.DealerClient, error) {
	params.SetDefaults()
	return s.Repo.GetDealerClients(ctx, params)
}



func (s *DealerClientService) GetAllDealerClients(ctx context.Context) ([]models.DealerClient, error) {
	return s.Repo.GetAll(ctx)
}

func (s *DealerClientService) UpdateDealerClient(ctx context.Context, id string, updates map[string]interface{}) error {
	return s.Repo.Update(ctx, id, updates)
}

func (s *DealerClientService) DeleteDealerClient(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

func (s *DealerClientService) GetDealerClientByID(ctx context.Context, id string) (models.DealerClient, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *DealerClientService) UpdateDealerClientStatus(ctx context.Context, id string, status string) error {
	return s.Repo.UpdateStatus(ctx, id, status)
}