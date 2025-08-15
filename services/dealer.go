package services

import (
	"context"
	"errors"
	"time"

	"myapp/models"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	 "go.mongodb.org/mongo-driver/bson"
)

type DealerService struct {
	DealerCollection *mongo.Collection
	JWTSecret      string
}

type Claims struct {
	ID    string `json:"id"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (d *DealerService) CreateDealer(ctx context.Context,dealer models.Dealer) error {
	

	hash, _ := bcrypt.GenerateFromPassword([]byte(dealer.Password), bcrypt.DefaultCost)
	dealer.Password = string(hash)

	_, err := d.DealerCollection.InsertOne(ctx, dealer)
	return err
}

func (d *DealerService) LoginDealer(ctx context.Context, phone, password string) (string, error) {
	

	var dbUser models.Dealer
	err := d.DealerCollection.FindOne(ctx, map[string]string{"phone": phone}).Decode(&dbUser)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	claims := &Claims{
		ID:    dbUser.ID.Hex(),
		Phone: dbUser.Phone,
		Role:  "dealer",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(d.JWTSecret))
}

func (d *DealerService) GetAllDealers(ctx context.Context) ([]models.Dealer, error) {
    cursor, err := d.DealerCollection.Find(ctx, bson.M{})
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
