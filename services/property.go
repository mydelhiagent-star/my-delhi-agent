package services

import (
	"context"
	"encoding/json"
	"fmt"
	"myapp/models"
	"myapp/repositories"
	"myapp/utils"
	"time"

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
	// Get property first to find dealer ID for cache invalidation
	property, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	// Update the property
	err = s.Repo.Update(context.Background(), id, updates)
	if err != nil {
		return err
	}

	// Invalidate dealer's property cache
	if s.RedisClient != nil {
		s.InvalidateDealerPropertyCache(property.DealerID)
	}

	return nil
}

func (s *PropertyService) DeleteProperty(id string) error {
	// Get property first to find dealer ID for cache invalidation
	property, err := s.Repo.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	// Delete the property
	err = s.Repo.Delete(context.Background(), id)
	if err != nil {
		return err
	}

	// Invalidate dealer's property cache
	if s.RedisClient != nil {
		s.InvalidateDealerPropertyCache(property.DealerID)
	}

	return nil
}



func (s *PropertyService) GetPropertiesByDealer(ctx context.Context, dealerID string, page, limit int) ([]models.Property, error) {
	// Validate inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 12 // Default limit
	}

	redisKey := fmt.Sprintf("properties_by_dealer:%s:page:%d:limit:%d", dealerID, page, limit)
	if s.RedisClient != nil {
		cached, err := s.RedisClient.Get(ctx, redisKey).Result()
		if err == nil {
			var properties []models.Property
			if json.Unmarshal([]byte(cached), &properties) == nil {
				return properties, nil
			}
		}
	}

	properties, err := s.Repo.GetByDealer(ctx, dealerID, page, limit)
	if err != nil {
		return nil, err
	}

	if s.RedisClient != nil {
		if data, err := json.Marshal(properties); err == nil {
			s.RedisClient.Set(ctx, redisKey, data, 30*time.Minute)
		}
	}

	return properties, nil
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