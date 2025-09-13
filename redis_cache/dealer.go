package redis_cache

import (
    "context"
    "fmt"
    "myapp/models"
    "time"
)

type DealerCache struct {
    cacheManager *CacheManager
}

func NewDealerCache(cacheManager *CacheManager) *DealerCache {
    return &DealerCache{
        cacheManager: cacheManager,
    }
}

func (dc *DealerCache) GetDealerByID(ctx context.Context, dealerID string) (models.Dealer, error) {
    key := fmt.Sprintf("dealer:id:%s", dealerID)
    
    var dealer models.Dealer
    err := dc.cacheManager.Get(ctx, key, &dealer)
    if err != nil {
        return models.Dealer{}, err // Cache miss
    }
    
    return dealer, nil
}

func (dc *DealerCache) SetDealerByID(ctx context.Context, dealerID string, dealer models.Dealer) error {
    key := fmt.Sprintf("dealer:id:%s", dealerID)
    return dc.cacheManager.Set(ctx, key, dealer, 30*time.Minute)
}

func (dc *DealerCache) GetDealerByPhone(ctx context.Context, phone string) (models.Dealer, error) {
    key := fmt.Sprintf("dealer:phone:%s", phone)
    
    var dealer models.Dealer
    err := dc.cacheManager.Get(ctx, key, &dealer)
    if err != nil {
        return models.Dealer{}, err // Cache miss
    }
    
    return dealer, nil
}

func (dc *DealerCache) SetDealerByPhone(ctx context.Context, phone string, dealer models.Dealer) error {
    key := fmt.Sprintf("dealer:phone:%s", phone)
    return dc.cacheManager.Set(ctx, key, dealer, 30*time.Minute)
}

func (dc *DealerCache) InvalidateDealer(ctx context.Context, dealerID string) error {
    // Invalidate all dealer-related cache entries
    patterns := []string{
        fmt.Sprintf("dealer:id:%s", dealerID),
        fmt.Sprintf("dealer:phone:*"), // This will invalidate all phone-based caches
    }
    
    for _, pattern := range patterns {
        if err := dc.cacheManager.DeleteByPattern(ctx, pattern); err != nil {
            return err
        }
    }
    
    return nil
}