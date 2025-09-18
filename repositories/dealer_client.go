package repositories

import (
	"context"
	"myapp/models"
)

type DealerClientRepository interface {
	Create(ctx context.Context, dealerClient models.DealerClient) (string, error)
	GetByID(ctx context.Context, id string) (models.DealerClient, error)
	GetDealerClients(ctx context.Context, params models.DealerClientQueryParams, fields []string) ([]models.DealerClient, error)
	GetAll(ctx context.Context) ([]models.DealerClient, error)
	Update(ctx context.Context, id string, updates models.DealerClientUpdate) error
	Delete(ctx context.Context, id string) error
	CheckPhoneExistsForDealer(ctx context.Context, dealerID, phone string) (bool, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	CheckPropertyInterestExists(ctx context.Context, dealerClientID string, propertyID string) (bool, error)
	CreateDealerClientPropertyInterest(ctx context.Context, dealerClientID string, dealerClientPropertyInterest models.DealerClientPropertyInterest) error
}
