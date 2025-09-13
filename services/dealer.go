package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"myapp/mongo_models"
	"myapp/repositories"
	"myapp/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type DealerService struct {
	DealerRepo repositories.DealerRepository
	TokenRepo  repositories.TokenRepository
	JWTSecret  string
}

func (s *DealerService) CreateDealer(ctx context.Context, dealer models.Dealer) error {
	// ← HASH password before insertion
	hash, err := utils.HashPassword(dealer.Password)
	if err != nil {
		return err
	}
	dealer.Password = string(hash)

	// ← Insert dealer (relying on MongoDB's unique constraints)
	_, err = s.DealerRepo.Create(ctx, dealer)
	if err != nil {
		// Check if the error is a unique constraint violation (e.g., phone or sublocation already exists)
		if mongo.IsDuplicateKeyError(err) {
			// Parse the error details to find out which field caused the conflict
			if strings.Contains(err.Error(), "phone") {
				return fmt.Errorf("phone number already exists")
			}
			if strings.Contains(err.Error(), "sub_location") {
				return fmt.Errorf("sublocation already exists")
			}
		}
		return fmt.Errorf("failed to create dealer")
	}

	return nil
}

func (s *DealerService) LoginDealer(ctx context.Context, phone, password string) (string, error) {

	dbUser, err := s.DealerRepo.GetByPhone(ctx, phone)
	if err != nil {
		return "", errors.New("invalid phone number")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	claims := &models.Claims{
		ID:    dbUser.ID.Hex(),
		Phone: dbUser.Phone,
		Role:  "dealer",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()), // unique timestamp
			ID:       uuid.New().String(),            // unique jti
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", err
	}
	tokenModel := models.Token{
		Token: tokenString,
		User:  dbUser.ID.Hex(),
	}
	err = s.TokenRepo.Create(ctx, tokenModel)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *DealerService) LogoutDealer(ctx context.Context, token string) error {
	return s.TokenRepo.Delete(ctx, token)
}

func (s *DealerService) GetAllDealers(ctx context.Context) ([]models.Dealer, error) {
	return s.DealerRepo.GetAll(ctx)
}

func (s *DealerService) GetDealersByLocation(ctx context.Context, subLocation string) ([]models.Dealer, error) {
	return s.DealerRepo.GetByLocation(ctx, subLocation)
}

func (s *DealerService) GetLocationsWithSubLocations(ctx context.Context) ([]models.LocationWithSubLocations, error) {
	return s.DealerRepo.GetLocationsWithSubLocations(ctx)
}

func (s *DealerService) DealerExists(ctx context.Context, dealerID primitive.ObjectID) (bool, error) {
	return s.DealerRepo.Exists(ctx, dealerID)
}

func (s *DealerService) GetDealerWithProperties(ctx context.Context, subLocation string) ([]map[string]interface{}, error) {
	// This would need a complex aggregation pipeline
	// For now, return empty slice
	return []map[string]interface{}{}, nil
}

func (s *DealerService) UpdateDealer(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return s.DealerRepo.Update(ctx, id, updates)
}

func (s *DealerService) DeleteDealer(ctx context.Context, id primitive.ObjectID) error {
	return s.DealerRepo.Delete(ctx, id)
}

func (s *DealerService) ResetPasswordDealer(ctx context.Context, dealerID primitive.ObjectID, newPassword string) error {
	// Hash the new password
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the dealer's password
	updates := map[string]interface{}{
		"password": string(hash),
	}
	return s.DealerRepo.Update(ctx, dealerID, updates)
}
