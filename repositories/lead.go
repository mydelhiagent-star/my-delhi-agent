package repositories

import (
	"context"
	"myapp/models"
)

type LeadRepository interface {
	Create(ctx context.Context, lead models.Lead) (string, error)
	GetByID(ctx context.Context, id string) (models.Lead, error)
	GetAll(ctx context.Context) ([]models.Lead, error)
	GetByDealerID(ctx context.Context, dealerID string) ([]models.Lead, error)
	GetLeads(ctx context.Context, params models.LeadQueryParams) ([]models.Lead, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	AddPropertyInterest(ctx context.Context, leadID string, propertyInterest models.PropertyInterest) error
	UpdatePropertyInterest(ctx context.Context, leadID, propertyID string, status, note string) error
	
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)
}
