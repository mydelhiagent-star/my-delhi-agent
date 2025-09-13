package mongo_repositories

import (
	"context"
	"myapp/models"
	mongoModels "myapp/mongo_models"
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

func (r *MongoDealerClientRepository) Create(ctx context.Context, dealerClient models.DealerClient) (string, error) {

	dealerObjectID, err := primitive.ObjectIDFromHex(dealerClient.DealerID)
	if err != nil {
		return "", err
	}

	propertyObjectID, err := primitive.ObjectIDFromHex(dealerClient.PropertyID)
	if err != nil {
		return "", err
	}

	mongoDealerClient := mongoModels.DealerClient{
		DealerID:   dealerObjectID,
		PropertyID: propertyObjectID,
		Name:       dealerClient.Name,
		Phone:      dealerClient.Phone,
		Note:       dealerClient.Note,
		Status:     dealerClient.Status,
	}

	result, err := r.dealerClientCollection.InsertOne(ctx, mongoDealerClient)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoDealerClientRepository) GetByID(ctx context.Context, id string) (models.DealerClient, error) {
	
	var dealerClient models.DealerClient
	err := r.dealerClientCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&dealerClient)
	if err != nil {
		return models.DealerClient{}, err
	}
	return dealerClient, nil
}

func (r *MongoDealerClientRepository) GetByDealerID(ctx context.Context, dealerID string) ([]models.DealerClient, error) {
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

func (r *MongoDealerClientRepository) GetByPropertyID(ctx context.Context, propertyID string) ([]models.DealerClient, error) {
	propertyObjectID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.dealerClientCollection.Find(ctx, bson.M{"property_id": propertyObjectID})
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

func (r *MongoDealerClientRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": updates}
	_, err = r.dealerClientCollection.UpdateByID(ctx, objectID, update)
	return err
}

func (r *MongoDealerClientRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.dealerClientCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoDealerClientRepository) CheckPhoneExistsForDealer(ctx context.Context, dealerID string, propertyID string, phone string) (bool, error) {
	dealerObjectID, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		return false, err
	}

	propertyObjectID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"dealer_id":   dealerObjectID,
		"property_id": propertyObjectID,
		"phone":       phone,
	}

	count, err := r.dealerClientCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *MongoDealerClientRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err = r.dealerClientCollection.UpdateByID(ctx, objectID, update)
	return err
}
