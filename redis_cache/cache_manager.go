package redis_cache



import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/go-redis/redis/v8"
)

type CacheManager struct {
    redisClient *redis.Client
}

func NewCacheManager(redisClient *redis.Client) *CacheManager {
    return &CacheManager{
        redisClient: redisClient,
    }
}

func (cm *CacheManager) Get(ctx context.Context, key string, result interface{}) error {
    if cm.redisClient == nil {
        return fmt.Errorf("redis client not available")
    }
    
    cached, err := cm.redisClient.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    
    return json.Unmarshal([]byte(cached), result)
}

func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    if cm.redisClient == nil {
        return fmt.Errorf("redis client not available")
    }
    
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return cm.redisClient.Set(ctx, key, data, ttl).Err()
}

func (cm *CacheManager) Delete(ctx context.Context, key string) error {
    if cm.redisClient == nil {
        return fmt.Errorf("redis client not available")
    }
    
    return cm.redisClient.Del(ctx, key).Err()
}

func (cm *CacheManager) DeleteByPattern(ctx context.Context, pattern string) error {
    if cm.redisClient == nil {
        return fmt.Errorf("redis client not available")
    }
    
    keys, err := cm.redisClient.Keys(ctx, pattern).Result()
    if err != nil {
        return err
    }
    
    if len(keys) > 0 {
        return cm.redisClient.Del(ctx, keys...).Err()
    }
    
    return nil
}

func (cm *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
    if cm.redisClient == nil {
        return false, fmt.Errorf("redis client not available")
    }
    
    count, err := cm.redisClient.Exists(ctx, key).Result()
    return count > 0, err
}