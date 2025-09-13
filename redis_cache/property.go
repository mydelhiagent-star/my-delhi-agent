package redis_cache

import (
    "context"
    "fmt"
    "myapp/models"
    "time"
)

type PropertyCache struct {
    cacheManager *CacheManager
}

func NewPropertyCache(cacheManager *CacheManager) *PropertyCache {
    return &PropertyCache{
        cacheManager: cacheManager,
    }
}

func (pc *PropertyCache) GetPropertiesByDealer(ctx context.Context, dealerID string, page, limit int) ([]models.Property, error) {
    key := fmt.Sprintf("properties:dealer:%s:page:%d:limit:%d", dealerID, page, limit)
    
    var properties []models.Property
    err := pc.cacheManager.Get(ctx, key, &properties)
    if err != nil {
        return nil, err // Cache miss
    }
    
    return properties, nil
}

func (pc *PropertyCache) SetPropertiesByDealer(ctx context.Context, dealerID string, page, limit int, properties []models.Property) error {
    key := fmt.Sprintf("properties:dealer:%s:page:%d:limit:%d", dealerID, page, limit)
    return pc.cacheManager.Set(ctx, key, properties, 30*time.Minute)
}

func (pc *PropertyCache) GetPropertyByID(ctx context.Context, propertyID string) (models.Property, error) {
    key := fmt.Sprintf("properties:id:%s", propertyID)
    
    var property models.Property
    err := pc.cacheManager.Get(ctx, key, &property)
    if err != nil {
        return models.Property{}, err // Cache miss
    }
    
    return property, nil
}

func (pc *PropertyCache) SetPropertyByID(ctx context.Context, propertyID string, property models.Property) error {
    key := fmt.Sprintf("properties:id:%s", propertyID)
    return pc.cacheManager.Set(ctx, key, property, 30*time.Minute)
}

func (pc *PropertyCache) InvalidateDealerProperties(ctx context.Context, dealerID string) error {
    pattern := fmt.Sprintf("properties:dealer:%s:*", dealerID)
    return pc.cacheManager.DeleteByPattern(ctx, pattern)
}

func (pc *PropertyCache) InvalidateProperty(ctx context.Context, propertyID string) error {
    key := fmt.Sprintf("properties:id:%s", propertyID)
    return pc.cacheManager.Delete(ctx, key)
}