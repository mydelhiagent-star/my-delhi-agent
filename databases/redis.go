package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis(uri string, username string, password string) (*redis.Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    client := redis.NewClient(&redis.Options{
        Addr:     uri,
        Username: username,
        Password: password,
        DB:       0,
    })

    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("Redis connection failed: %w", err)
    }

    return client, nil
}