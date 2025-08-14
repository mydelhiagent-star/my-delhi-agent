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

type AuthService struct {
	UserCollection *mongo.Collection
	JWTSecret      string
}

type Claims struct {
	ID    string `json:"id"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) CreateUser(ctx context.Context,user models.User) error {
	

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	_, err := s.UserCollection.InsertOne(ctx, user)
	return err
}

func (s *AuthService) Login(phone, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbUser models.User
	err := s.UserCollection.FindOne(ctx, map[string]string{"phone": phone}).Decode(&dbUser)
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
		Role:  dbUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}

func (s *AuthService) GetAllUsers(ctx context.Context) ([]models.User, error) {
    cursor, err := s.UserCollection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var users []models.User
    for cursor.Next(ctx) {
        var user models.User
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
