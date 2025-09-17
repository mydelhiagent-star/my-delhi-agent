package repositories

import (
	"context"
	"myapp/models"
)

type PropertyRepository interface {
	Create(ctx context.Context, property models.Property) (string, error)
	GetByID(ctx context.Context, id string) (models.Property, error)
	GetByDealer(ctx context.Context, dealerID string, page, limit int) ([]models.Property, error)
	Update(ctx context.Context, id string, updates models.PropertyUpdate) error
	Delete(ctx context.Context, id string) error
	GetNextPropertyNumber(ctx context.Context) (int64, error)
	GetProperties(ctx context.Context, params models.PropertyQueryParams, fields []string) ([]models.Property, error)
}
