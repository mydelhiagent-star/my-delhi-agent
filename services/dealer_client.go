package services

import (
	"context"
	"errors"
	"myapp/models"
	"myapp/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DealerClientService struct {
	Repo repositories.DealerClientRepository
	PropertyRepo repositories.PropertyRepository
}

func (s *DealerClientService) CheckPhoneExistsForDealer(ctx context.Context, dealerID string, phone string) (bool, error) {
	return s.Repo.CheckPhoneExistsForDealer(ctx, dealerID, phone)
}

func (s *DealerClientService) CreateDealerClient(ctx context.Context, dealerClient models.DealerClient) (string, error) {
	
	exists, err := s.CheckPhoneExistsForDealer(ctx, dealerClient.DealerID, dealerClient.Phone)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("phone number already exists")
	}

	

	return s.Repo.Create(ctx, dealerClient)
}

func (s *DealerClientService) GetDealerClients(ctx context.Context, params models.DealerClientQueryParams, fields []string) ([]models.DealerClient, error) {
	params.SetDefaults()

	dealerClients, err := s.Repo.GetDealerClients(ctx, params, fields)
	if err != nil {
		return nil, err
	}

	if len(dealerClients) == 0 {
		return dealerClients, nil
	}

	
	propertyIDs := make([]primitive.ObjectID, 0)
	for _, dealerClient := range dealerClients {
		for _, propertyInterest := range dealerClient.PropertyInterests {
			if objectID, err := primitive.ObjectIDFromHex(propertyInterest.PropertyID); err == nil {
				propertyIDs = append(propertyIDs, objectID)
			}
		}
	}
	
	if len(propertyIDs) == 0 {
		return dealerClients, nil
	}

	// Find deleted/sold properties
	filter := bson.M{
		"_id": bson.M{"$in": propertyIDs},
		 "$or": bson.A{
			bson.M{"is_deleted": true},
			bson.M{"sold": true},
		},
	}

	projection := bson.M{
		"_id": 1,
	}

	properties, err := s.PropertyRepo.GetFilteredProperties(ctx, filter, projection, int64(len(propertyIDs)), 0)
	if err != nil {
		return nil, err
	}

	
	propertiesToRemove := make(map[string]bool)
	for _, property := range properties {
		propertiesToRemove[property.ID] = true
	}

	
	for i := range dealerClients {
		var filteredInterests []models.DealerClientPropertyInterest
		
		for _, interest := range dealerClients[i].PropertyInterests {
			
			if !propertiesToRemove[interest.PropertyID] {
				filteredInterests = append(filteredInterests, interest)
			}
		}
		
		dealerClients[i].PropertyInterests = filteredInterests
	}

	return dealerClients, nil
}





func (s *DealerClientService) UpdateDealerClient(ctx context.Context, id string, updates models.DealerClientUpdate) error {
	return s.Repo.Update(ctx, id, updates)
}

func (s *DealerClientService) DeleteDealerClient(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

func (s *DealerClientService) GetDealerClientByID(ctx context.Context, id string) (models.DealerClient, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *DealerClientService) UpdateDealerClientStatus(ctx context.Context, id string, status string) error {
	return s.Repo.UpdateStatus(ctx, id, status)
}

func (s *DealerClientService) CreateDealerClientPropertyInterest(ctx context.Context, dealerClientID string, dealerClientPropertyInterest models.DealerClientPropertyInterest) error {
	
	exists, err := s.Repo.CheckPropertyInterestExists(ctx, dealerClientID, dealerClientPropertyInterest.PropertyID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("client is already added to this property")
	}

	return s.Repo.CreateDealerClientPropertyInterest(ctx, dealerClientID, dealerClientPropertyInterest)
}

func (s *DealerClientService) UpdateDealerClientPropertyInterest(ctx context.Context, dealerClientID string, propertyInterestID string, update models.DealerClientPropertyInterestUpdate) error {
	return s.Repo.UpdateDealerClientPropertyInterest(ctx, dealerClientID, propertyInterestID, update)
}

func (s *DealerClientService) DeleteDealerClientPropertyInterest(ctx context.Context, dealerClientID string, propertyInterestID string) error {
	return s.Repo.DeleteDealerClientPropertyInterest(ctx, dealerClientID, propertyInterestID)
}