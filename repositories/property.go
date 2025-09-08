package repositories

import (
	"context"
	"myapp/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PropertyRepository defines the interface for property data operations
type PropertyRepository interface {
	Create(ctx context.Context, property models.Property) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Property, error)
	GetByNumber(ctx context.Context, propertyNumber int64) (*models.Property, error)
	GetByDealer(ctx context.Context, dealerID primitive.ObjectID, page, limit int) ([]models.Property, error)
	GetAll(ctx context.Context) ([]models.Property, error)
	Search(ctx context.Context, filter bson.M, page, limit int, fields []string) ([]models.Property, error)
	Update(ctx context.Context, id primitive.ObjectID, updates models.PropertyUpdate) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetNextPropertyNumber(ctx context.Context) (int64, error)
		
}
