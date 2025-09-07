package repositories

import (
	"context"
	"myapp/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DealerClientRepository defines the interface for dealer client data operations
type DealerClientRepository interface {
	Create(ctx context.Context, dealerClient models.DealerClient) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.DealerClient, error)
	GetByDealerID(ctx context.Context, dealerID primitive.ObjectID) ([]models.DealerClient, error)
	GetByPropertyID(ctx context.Context, propertyID primitive.ObjectID) ([]models.DealerClient, error)
	GetAll(ctx context.Context) ([]models.DealerClient, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
