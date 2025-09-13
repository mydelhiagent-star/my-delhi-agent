package mongo_repositories

// import (
// 	"context"
// 	"myapp/models"
// 	mongoModels "myapp/mongo_models"
// 	"myapp/repositories"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type MongoTokenRepository struct {
// 	tokenCollection *mongo.Collection
// }

// func NewMongoTokenRepository(tokenCollection *mongo.Collection) repositories.TokenRepository {
// 	return &MongoTokenRepository{
// 		tokenCollection: tokenCollection,
// 	}
// }

// func (r *MongoTokenRepository) Create(ctx context.Context, token models.Token) error {
// 	_, err := r.tokenCollection.InsertOne(ctx, token)
// 	return err
// }

// func (r *MongoTokenRepository) GetByToken(ctx context.Context, token string) (*models.Token, error) {
// 	var tokenModel models.Token
// 	err := r.tokenCollection.FindOne(ctx, bson.M{"token": token}).Decode(&tokenModel)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &tokenModel, nil
// }

// func (r *MongoTokenRepository) Delete(ctx context.Context, token string) error {
// 	_, err := r.tokenCollection.DeleteOne(ctx, bson.M{"token": token})
// 	return err
// }

// func (r *MongoTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
// 	_, err := r.tokenCollection.DeleteMany(ctx, bson.M{"user_id": userID})
// 	return err
// }
