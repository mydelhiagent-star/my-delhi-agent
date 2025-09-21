package services

import (
	"context"
	"myapp/models"
	"myapp/repositories"
)

type InquiryService struct {
	inquiryRepo repositories.InquiryRepository
}

func NewInquiryService(inquiryRepo repositories.InquiryRepository) *InquiryService {
	return &InquiryService{
		inquiryRepo: inquiryRepo,
	}
}

func (s *InquiryService) CreateInquiry(ctx context.Context, inquiry models.Inquiry) (models.Inquiry, error) {
	return s.inquiryRepo.Create(ctx, inquiry)
}

func (s *InquiryService) GetInquiryByID(ctx context.Context, id string) (models.Inquiry, error) {
	return s.inquiryRepo.GetByID(ctx, id)
}

func (s *InquiryService) GetAllInquiries(ctx context.Context, params models.InquiryQueryParams) ([]models.Inquiry, error) {
	return s.inquiryRepo.GetAll(ctx, params)
}

func (s *InquiryService) UpdateInquiry(ctx context.Context, id string, updates models.InquiryUpdate) error {
	return s.inquiryRepo.Update(ctx, id, updates)
}

func (s *InquiryService) DeleteInquiry(ctx context.Context, id string) error {
	return s.inquiryRepo.Delete(ctx, id)
}
