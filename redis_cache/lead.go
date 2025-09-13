package redis_cache

import (
    "context"
    "fmt"
    "myapp/models"
    "time"
)

type LeadCache struct {
    cacheManager *CacheManager
}

func NewLeadCache(cacheManager *CacheManager) *LeadCache {
    return &LeadCache{
        cacheManager: cacheManager,
    }
}

func (lc *LeadCache) GetLeadByID(ctx context.Context, leadID string) (models.Lead, error) {
    key := fmt.Sprintf("lead:id:%s", leadID)
    
    var lead models.Lead
    err := lc.cacheManager.Get(ctx, key, &lead)
    if err != nil {
        return models.Lead{}, err // Cache miss
    }
    
    return lead, nil
}

func (lc *LeadCache) SetLeadByID(ctx context.Context, leadID string, lead models.Lead) error {
    key := fmt.Sprintf("lead:id:%s", leadID)
    return lc.cacheManager.Set(ctx, key, lead, 30*time.Minute)
}

func (lc *LeadCache) GetLeadsByDealer(ctx context.Context, dealerID string, page, limit int) ([]models.Lead, error) {
    key := fmt.Sprintf("leads:dealer:%s:page:%d:limit:%d", dealerID, page, limit)
    
    var leads []models.Lead
    err := lc.cacheManager.Get(ctx, key, &leads)
    if err != nil {
        return nil, err // Cache miss
    }
    
    return leads, nil
}

func (lc *LeadCache) SetLeadsByDealer(ctx context.Context, dealerID string, page, limit int, leads []models.Lead) error {
    key := fmt.Sprintf("leads:dealer:%s:page:%d:limit:%d", dealerID, page, limit)
    return lc.cacheManager.Set(ctx, key, leads, 30*time.Minute)
}

func (lc *LeadCache) InvalidateLead(ctx context.Context, leadID string) error {
    key := fmt.Sprintf("lead:id:%s", leadID)
    return lc.cacheManager.Delete(ctx, key)
}

func (lc *LeadCache) InvalidateDealerLeads(ctx context.Context, dealerID string) error {
    pattern := fmt.Sprintf("leads:dealer:%s:*", dealerID)
    return lc.cacheManager.DeleteByPattern(ctx, pattern)
}