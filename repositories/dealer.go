package repositories

import (
	"context"
	"myapp/mongo_models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DealerRepository defines the interface for dealer data operations
type DealerRepository interface {
	Create(ctx context.Context, dealer models.Dealer) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Dealer, error)
	GetByPhone(ctx context.Context, phone string) (*models.Dealer, error)
	GetByEmail(ctx context.Context, email string) (*models.Dealer, error)
	GetAll(ctx context.Context) ([]models.Dealer, error)
	GetByLocation(ctx context.Context, subLocation string) ([]models.Dealer, error)
	GetLocationsWithSubLocations(ctx context.Context) ([]models.LocationWithSubLocations, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	Exists(ctx context.Context, id primitive.ObjectID) (bool, error)
}
