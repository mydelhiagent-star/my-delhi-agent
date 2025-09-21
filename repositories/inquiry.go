package repositories

import (
	"context"
	"myapp/models"
)

type InquiryRepository interface {
	Create(ctx context.Context, inquiry models.Inquiry) (models.Inquiry, error)
	GetByID(ctx context.Context, id string) (models.Inquiry, error)
	GetAll(ctx context.Context, params models.InquiryQueryParams) ([]models.Inquiry, error)
	Update(ctx context.Context, id string, updates models.InquiryUpdate) error
	Delete(ctx context.Context, id string) error
}
