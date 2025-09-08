package mongo_repository

import (
	"context"
	"myapp/models"
	"myapp/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDealerClientRepository struct {
	dealerClientCollection *mongo.Collection
}

func NewMongoDealerClientRepository(dealerClientCollection *mongo.Collection) repositories.DealerClientRepository {
	return &MongoDealerClientRepository{
		dealerClientCollection: dealerClientCollection,
	}
}

func (r *MongoDealerClientRepository) Create(ctx context.Context, dealerClient models.DealerClient) (primitive.ObjectID, error) {
	result, err := r.dealerClientCollection.InsertOne(ctx, dealerClient)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *MongoDealerClientRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.DealerClient, error) {
	var dealerClient models.DealerClient
	err := r.dealerClientCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&dealerClient)
	if err != nil {
		return nil, err
	}
	return &dealerClient, nil
}

func (r *MongoDealerClientRepository) GetByDealerID(ctx context.Context, dealerID primitive.ObjectID) ([]models.DealerClient, error) {
	cursor, err := r.dealerClientCollection.Find(ctx, bson.M{"dealer_id": dealerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dealerClients []models.DealerClient
	if err := cursor.All(ctx, &dealerClients); err != nil {
		return nil, err
	}
	return dealerClients, nil
}

func (r *MongoDealerClientRepository) GetByPropertyID(ctx context.Context, propertyID primitive.ObjectID) ([]models.DealerClient, error) {
	cursor, err := r.dealerClientCollection.Find(ctx, bson.M{"property_id": propertyID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dealerClients []models.DealerClient
	if err := cursor.All(ctx, &dealerClients); err != nil {
		return nil, err
	}
	return dealerClients, nil
}

func (r *MongoDealerClientRepository) GetAll(ctx context.Context) ([]models.DealerClient, error) {
	cursor, err := r.dealerClientCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dealerClients []models.DealerClient
	if err := cursor.All(ctx, &dealerClients); err != nil {
		return nil, err
	}
	return dealerClients, nil
}

func (r *MongoDealerClientRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	update := bson.M{"$set": updates}
	_, err := r.dealerClientCollection.UpdateByID(ctx, id, update)
	return err
}

func (r *MongoDealerClientRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.dealerClientCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoDealerClientRepository) CheckPhoneExistsForDealer(ctx context.Context, dealerID primitive.ObjectID, propertyID primitive.ObjectID, phone string) (bool, error) {
	filter := bson.M{
		"dealer_id":   dealerID,
		"property_id": propertyID,
		"phone":       phone,
	}

	count, err := r.dealerClientCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MongoDealerClientRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.dealerClientCollection.UpdateByID(ctx, id, update)
	return err
}