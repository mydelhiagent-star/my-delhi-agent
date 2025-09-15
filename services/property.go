package services

import (
	"context"
	
	"fmt"
	"myapp/models"
	"myapp/repositories"
	"myapp/utils"
	

	"github.com/go-redis/redis/v8"
)

type PropertyService struct {
	Repo        repositories.PropertyRepository
	RedisClient *redis.Client
}

func (s *PropertyService) CreateProperty(ctx context.Context, property models.Property) (string, error) {
	var resultID string

	err := utils.Retry(ctx, func() error {
		propertyNumber, err := s.Repo.GetNextPropertyNumber(ctx)
		if err != nil {
			return err
		}

		property.PropertyNumber = propertyNumber

		resultID, err = s.Repo.Create(ctx, property)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if s.RedisClient != nil {
		s.InvalidateDealerPropertyCache(property.DealerID)
	}

	return resultID, nil
}


func (s *PropertyService) UpdateProperty(id string, updates models.PropertyUpdate) error {
	
	property, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	
	err = s.Repo.Update(context.Background(), id, updates)
	if err != nil {
		return err
	}

	if s.RedisClient != nil {
		s.InvalidateDealerPropertyCache(property.DealerID)
	}

	return nil
}

func (s *PropertyService) DeleteProperty(id string) error {
	property, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	err = s.Repo.Delete(context.Background(), id)
	if err != nil {
		return err
	}

	if s.RedisClient != nil {
		s.InvalidateDealerPropertyCache(property.DealerID)
	}

	return nil
}







func (s *PropertyService) InvalidateDealerPropertyCache(dealerID string) {
	ctx := context.Background()

	// Invalidate all dealer property pages
	pattern := fmt.Sprintf("properties_by_dealer:%s:page:*", dealerID)
	keys, err := s.RedisClient.Keys(ctx, pattern).Result()
	if err == nil {
		for _, key := range keys {
			s.RedisClient.Del(ctx, key)
		}
	}
}

func (s *PropertyService) GetProperties(ctx context.Context, filters map[string]interface{}, page, limit int) ([]models.Property, error) {
	return s.Repo.GetProperties(ctx, filters, page, limit)
}