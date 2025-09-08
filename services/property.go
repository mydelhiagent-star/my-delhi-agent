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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PropertyService struct {
	Repo        repositories.PropertyRepository
	RedisClient *redis.Client
}

func (s *PropertyService) CreateProperty(ctx context.Context, property models.Property) (primitive.ObjectID, error) {
	var resultID primitive.ObjectID

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
		return primitive.NilObjectID, err
	}

	if s.RedisClient != nil {
		s.InvalidateDealerPropertyCache(property.DealerID)
	}

	return resultID, nil
}

func (s *PropertyService) GetPropertyByNumber(ctx context.Context, propertyNumber int64) (*models.Property, error) {
	return s.Repo.GetByNumber(ctx, propertyNumber)
}

func (s *PropertyService) GetPropertyByID(id primitive.ObjectID) (*models.Property, error) {
	return s.Repo.GetByID(context.Background(), id)
}

func (s *PropertyService) UpdateProperty(id primitive.ObjectID, updates models.PropertyUpdate) error {
	return s.Repo.Update(context.Background(), id, updates)
}

func (s *PropertyService) DeleteProperty(id primitive.ObjectID) error {
	return s.Repo.Delete(context.Background(), id)
}

func (s *PropertyService) GetAllProperties(ctx context.Context) ([]models.Property, error) {
	return s.Repo.GetAll(ctx)
}

func (s *PropertyService) GetPropertiesByDealer(ctx context.Context, dealerID primitive.ObjectID, page, limit int) ([]models.Property, error) {
	// ← VALIDATE inputs
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 12 // Default limit
	}

	redisKey := fmt.Sprintf("properties_by_dealer:%s:page:%d:limit:%d", dealerID.Hex(), page, limit)
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

func (s *PropertyService) SearchProperties(ctx context.Context, filter bson.M, page, limit int, fields []string) ([]models.Property, error) {
	return s.Repo.Search(ctx, filter, page, limit, fields)
}

func (s *PropertyService) InvalidateDealerPropertyCache(dealerID primitive.ObjectID) {
	ctx := context.Background()

	// Invalidate all dealer property pages
	pattern := fmt.Sprintf("properties_by_dealer:%s:page:*", dealerID.Hex())
	keys, err := s.RedisClient.Keys(ctx, pattern).Result()
	if err == nil {
		for _, key := range keys {
			s.RedisClient.Del(ctx, key)
		}
	}
}

func (s *PropertyService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Property, error) {
	return s.Repo.GetByID(ctx, id)
}
