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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type DealerService struct {
	DealerCollection *mongo.Collection
	PropertyCollection *mongo.Collection
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

// Replace in services/dealer.go
func (s *DealerService) GetDealerWithProperties(ctx context.Context, subLocation string) ([]map[string]interface{}, error) {
    // Step 1: Get dealer first (lightweight query)
    var dealer models.Dealer
    err := s.DealerCollection.FindOne(ctx, bson.M{"sub_location": subLocation}).Decode(&dealer)
    if err != nil {
        return nil, err
    }

    // Step 2: Get properties separately with pagination (prevents memory explosion)
    filter := bson.M{
        "dealer_id":  dealer.ID,
        "is_deleted": bson.M{"$ne": true},
        "sold":       bson.M{"$ne": true},
    }
    
    opts := options.Find().
        SetSort(bson.M{"_id": -1}).
        SetLimit(50).                    // Limit properties per dealer
        SetBatchSize(100)                // Batch size for large results
    
    cursor, err := s.PropertyCollection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var properties []models.Property
    if err := cursor.All(ctx, &properties); err != nil{
        return nil, err
    }

    // Step 3: Combine results
    result := map[string]interface{}{
        "dealer":     dealer,
        "properties": properties,
    }

    return []map[string]interface{}{result}, nil
}

func (s *DealerService) UpdateDealer(ctx context.Context, dealerID primitive.ObjectID, dealer models.Dealer) error {
	_, err := s.DealerCollection.UpdateByID(ctx, dealerID, bson.M{"$set": dealer})
	return err
}

func (s *DealerService) DeleteDealer(ctx context.Context, dealerID primitive.ObjectID) error {
	// ← START session for transaction
	session, err := s.DealerCollection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// ← EXECUTE transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Step 1: Mark all dealer's properties as deleted (soft delete)
		propertyCollection := s.DealerCollection.Database().Collection("property")
		propertyUpdate := bson.M{
			"$set": bson.M{
				"is_deleted": true,
			},
		}
		_, err := propertyCollection.UpdateMany(sessCtx,
			bson.M{"dealer_id": dealerID},
			propertyUpdate)
		if err != nil {
			return nil, err
		}

		// Step 2: Remove dealer's properties from leads' property interests
		leadCollection := s.DealerCollection.Database().Collection("leads")
		leadUpdate := bson.M{
			"$pull": bson.M{
				"properties": bson.M{
					"dealer_id": dealerID,
				},
			},
		}
		_, err = leadCollection.UpdateMany(sessCtx,
			bson.M{"properties.dealer_id": dealerID},
			leadUpdate)
		if err != nil {
			return nil, err
		}

		// Step 3: Delete the dealer
		_, err = s.DealerCollection.DeleteOne(sessCtx, bson.M{"_id": dealerID})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (s *DealerService) ResetPasswordDealer(ctx context.Context, dealerID primitive.ObjectID, newPassword string) error {
	// ← VALIDATE dealer exists
	var dealer models.Dealer
	err := s.DealerCollection.FindOne(ctx, bson.M{"_id": dealerID}).Decode(&dealer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("dealer not found")
		}
		return err
	}

	// ← HASH the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// ← UPDATE dealer's password
	update := bson.M{
		"$set": bson.M{
			"password": string(hashedPassword),
		},
	}

	_, err = s.DealerCollection.UpdateByID(ctx, dealerID, update)
	return err
}
