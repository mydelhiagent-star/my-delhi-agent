package services

import (
	"context"
	"errors"
	"time"

	"myapp/models"
	"myapp/repositories"

	"go.mongodb.org/mongo-driver/bson"
)

type LeadService struct {
	Repo repositories.LeadRepository
	PropertyRepo repositories.PropertyRepository
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

func (s *LeadService) GetLeads(ctx context.Context, params models.LeadQueryParams) ([]models.Lead, error) {
	leads, err := s.Repo.GetLeads(ctx, params)
	if err != nil {
		return nil, err
	}

	propertyIDs := make([]primitive.ObjectID, 0)
	for _, lead := range leads {
		for _, propertyInterest := range lead.Properties {
			propertyIDs = append(propertyIDs, propertyInterest.PropertyID)
		}
	}
	if len(propertyIDs) == 0 {
		return leads, nil
	}

	
	filter := bson.M{
		"_id": bson.M{"$in": propertyIDs},
		"$or": bson.A{
			bson.M{"is_deleted": true},
			bson.M{"sold": true},
		},
	}

	projection := bson.M{"_id": 1}
	properties, err := s.PropertyRepo.GetFilteredProperties(ctx, filter, projection, int64(len(objectIDs)), 0)
	if err != nil {
		return nil, err
	}

	propertiesToRemove := make(map[string]bool)
	for _, property := range properties {
		propertiesToRemove[property.ID] = true
	}

	// Filter out deleted/sold property interests
	for i := range leads {
		var filteredProperties []models.PropertyInterest
		
		for _, propertyInterest := range leads[i].Properties {
			if !propertiesToRemove[propertyInterest.PropertyID] {
				filteredProperties = append(filteredProperties, propertyInterest)
			}
		}
		
		leads[i].Properties = filteredProperties
	}

	return leads, nil

	
}


func (s *LeadService) UpdatePropertyInterest(ctx context.Context, leadID string, propertyID string, status string, note string) error {
	return s.Repo.UpdatePropertyInterest(ctx, leadID, propertyID, status, note)
}

