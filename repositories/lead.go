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
	Search(ctx context.Context, filter map[string]interface{}, page, limit int, fields []string) ([]models.Lead, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	AddPropertyInterest(ctx context.Context, leadID string, propertyInterest models.PropertyInterest) error
	UpdatePropertyInterest(ctx context.Context, leadID, propertyID string, status, note string) error
	GetLeadPropertyDetails(ctx context.Context, leadID string) ([]map[string]interface{}, error)
	GetDealerLeads(ctx context.Context, dealerID string) ([]models.Lead, error)
	GetPropertyDetails(ctx context.Context, soldStr, deletedStr string) ([]map[string]interface{}, error)
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)
}
