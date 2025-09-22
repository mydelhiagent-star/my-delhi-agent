package repositories

import (
	"context"
	"myapp/models"
)


type DealerRepository interface {
	Create(ctx context.Context, dealer models.Dealer) (string, error)
	GetByID(ctx context.Context, id string) (models.Dealer, error)
	GetByPhone(ctx context.Context, phone string) (models.Dealer, error)
	GetByEmail(ctx context.Context, email string) (models.Dealer, error)
	GetAll(ctx context.Context) ([]models.Dealer, error)
	GetByLocation(ctx context.Context, subLocation string) ([]models.Dealer, error)
	GetLocationsWithSubLocations(ctx context.Context) ([]models.LocationWithSubLocations, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	GetDealerWithProperties(ctx context.Context, subLocation string) ([]map[string]interface{}, error)
}
