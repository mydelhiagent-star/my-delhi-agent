package repositories

import (
	"context"
	"myapp/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LeadRepository defines the interface for lead data operations
type LeadRepository interface {
	Create(ctx context.Context, lead models.Lead) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (models.Lead, error)
	GetAll(ctx context.Context) ([]models.Lead, error)
	GetByDealerID(ctx context.Context, dealerID primitive.ObjectID) ([]models.Lead, error)
	Search(ctx context.Context, filter bson.M, page, limit int, fields []string) ([]models.Lead, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	AddPropertyInterest(ctx context.Context, leadID primitive.ObjectID, propertyInterest models.PropertyInterest) error
	UpdatePropertyInterest(ctx context.Context, leadID, propertyID primitive.ObjectID, status, note string) error
	GetLeadPropertyDetails(ctx context.Context, leadID primitive.ObjectID) ([]bson.M, error)
	GetDealerLeads(ctx context.Context, dealerID primitive.ObjectID) ([]models.Lead, error)
	GetPropertyDetails(ctx context.Context, soldStr, deletedStr string) ([]bson.M, error)
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)
}
