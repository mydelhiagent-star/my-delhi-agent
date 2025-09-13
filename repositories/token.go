package repositories

import (
	"context"
	"myapp/mongo_models"
)

// TokenRepository defines the interface for token data operations
type TokenRepository interface {
	Create(ctx context.Context, token models.Token) error
	GetByToken(ctx context.Context, token string) (*models.Token, error)
	Delete(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
