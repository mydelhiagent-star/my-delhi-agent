package services

import (
	"context"
	"errors"
	"time"

	"myapp/constants"
	"myapp/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type DealerService struct {
	DealerCollection *mongo.Collection
	TokenCollection  *mongo.Collection
	JWTSecret        string
}

func (s *DealerService) CreateDealer(ctx context.Context, dealer models.Dealer) error {
	if !constants.IsValidLocation(dealer.Location) {
		return errors.New("invalid location")
	}
	existingDealer := models.Dealer{}
	err := s.DealerCollection.FindOne(ctx, bson.M{"phone": dealer.Phone}).Decode(&existingDealer)
	if err == nil {
		return errors.New("phone number already exists")
	}
	if err != mongo.ErrNoDocuments {
		return errors.New("database error checking phone number" + err.Error())
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(dealer.Password), bcrypt.DefaultCost)
	dealer.Password = string(hash)

	_, err = s.DealerCollection.InsertOne(ctx, dealer)
	return err
}

func (s *DealerService) LoginDealer(ctx context.Context, phone, password string) (string, error) {

	var dbUser models.Dealer
	err := s.DealerCollection.FindOne(ctx, map[string]string{"phone": phone}).Decode(&dbUser)
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
	_, err = s.TokenCollection.InsertOne(ctx, bson.M{
		"token": tokenString,
		"user":  dbUser.ID.Hex(),
	})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *DealerService) LogoutDealer(ctx context.Context, token string) error {
	_, err := s.TokenCollection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

func (s *DealerService) GetAllDealers(ctx context.Context) ([]models.Dealer, error) {
	cursor, err := s.DealerCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.Dealer
	for cursor.Next(ctx) {
		var user models.Dealer
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *DealerService) GetDealersByLocation(ctx context.Context, subLocation string) ([]models.Dealer, error) {

	filter := bson.M{
		"sub_location": subLocation,
	}

	cursor, err := s.DealerCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dealers []models.Dealer
	if err := cursor.All(ctx, &dealers); err != nil {
		return nil, err
	}

	return dealers, nil
}

func (s *DealerService) GetLocationsWithSubLocations(ctx context.Context) ([]models.LocationWithSubLocations, error) {
	// Fetch all dealers with location present in constants.Locations
	cursor, err := s.DealerCollection.Find(
		ctx,
		bson.M{"location": bson.M{"$in": constants.Locations}},
		options.Find().SetProjection(bson.M{
			"location":     1,
			"sub_location": 1,
			"_id":          0, // exclude _id if not needed
		}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Map: location -> set of sublocations
	locationMap := make(map[string]map[string]struct{})

	for cursor.Next(ctx) {
		var dealer struct {
			Location    string `bson:"location"`
			SubLocation string `bson:"sub_location"`
		}
		if err := cursor.Decode(&dealer); err != nil {
			return nil, err
		}

		if dealer.Location != "" && dealer.SubLocation != "" {
			if _, ok := locationMap[dealer.Location]; !ok {
				locationMap[dealer.Location] = make(map[string]struct{})
			}
			locationMap[dealer.Location][dealer.SubLocation] = struct{}{}
		}
	}

	// Prepare result
	result := make([]models.LocationWithSubLocations, 0, len(constants.Locations))
	for _, loc := range constants.Locations {
		subLocs := make([]string, 0)
		if subs, ok := locationMap[loc]; ok {
			for sub := range subs {
				subLocs = append(subLocs, sub)
			}
		}
		result = append(result, models.LocationWithSubLocations{
			Location:    loc,
			SubLocation: subLocs,
		})
	}

	return result, nil
}

// In services/dealer_service.go

func (s *DealerService) GetDealerWithProperties(ctx context.Context, subLocation string) ([]map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"sub_location": subLocation}}},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "property",
				"localField":   "_id",
				"foreignField": "dealer_id",
				"as":           "properties",
			}},
		},
	}

	cursor, err := s.DealerCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil // one dealer per subLocation
}
