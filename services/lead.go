package services

import (
	"context"
	"errors"
	"time"

	"myapp/models"
	"myapp/repositories"
)

type LeadService struct {
	Repo repositories.LeadRepository
}

func (s *LeadService) CreateLead(ctx context.Context, lead models.Lead) (string, error) {
	exists, err := s.Repo.CheckPhoneExists(ctx, lead.Phone)
	if err != nil {
		return "", errors.New("database error checking phone")
	}
	if exists {
		return "", errors.New("lead with phone already exists")
	}

	return s.Repo.Create(ctx, lead)
}

func (s *LeadService) GetLeadByID(ctx context.Context, id string) (models.Lead, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *LeadService) GetAllLeads(ctx context.Context) ([]models.Lead, error) {
	return s.Repo.GetAll(ctx)
}

func (s *LeadService) GetAllLeadsByDealerID(ctx context.Context, dealerID string) ([]models.Lead, error) {
	return s.Repo.GetByDealerID(ctx, dealerID)
}

func (s *LeadService) UpdateLead(ctx context.Context, id string, updateData map[string]interface{}) error {
	return s.Repo.Update(ctx, id, updateData)
}

func (s *LeadService) DeleteLead(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

func (s *LeadService) AddPropertyInterest(ctx context.Context, leadID string, propertyInterest models.PropertyInterest) error {
	// Set status
	propertyInterest.Status = "view"
	propertyInterest.CreatedAt = time.Now()

	return s.Repo.AddPropertyInterest(ctx, leadID, propertyInterest)
}

func (s *LeadService) SearchLeads(ctx context.Context, filter map[string]interface{}, page, limit int, fields []string) ([]models.Lead, error) {
	// Validate inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return s.Repo.Search(ctx, filter, page, limit, fields)
}

func (s *LeadService) GetLeadPropertyDetails(ctx context.Context, leadID string) ([]map[string]interface{}, error) {
	return s.Repo.GetLeadPropertyDetails(ctx, leadID)
}

func (s *LeadService) GetPropertyDetails(ctx context.Context, soldStr string, deletedStr string) ([]map[string]interface{}, error) {
	return s.Repo.GetPropertyDetails(ctx, soldStr, deletedStr)
}

func (s *LeadService) UpdatePropertyInterest(ctx context.Context, leadID string, propertyID string, status string, note string) error {
	return s.Repo.UpdatePropertyInterest(ctx, leadID, propertyID, status, note)
}

func (s *LeadService) GetDealerLeads(ctx context.Context, dealerID string) ([]models.Lead, error) {
	return s.Repo.GetDealerLeads(ctx, dealerID)
}