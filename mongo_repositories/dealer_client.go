package mongo_repositories

import (
	"context"
	"myapp/converters"
	"myapp/models"
	mongoModels "myapp/mongo_models"
	"myapp/repositories"
	"myapp/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	mongoPropertyInterests := converters.ToMongoDealerClientPropertyInterestSlice(dealerClient.PropertyInterests)

	

	mongoDealerClient := mongoModels.DealerClient{
		DealerID:   dealerObjectID,
		Name:       dealerClient.Name,
		Phone:      dealerClient.Phone,
		Note:       dealerClient.Note,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		PropertyInterests: mongoPropertyInterests,
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

func (r *MongoDealerClientRepository) GetDealerClients(ctx context.Context, params models.DealerClientQueryParams, fields []string) ([]models.DealerClient, error) {
	
	filter := utils.BuildMongoFilter(params)


	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(*params.Page - 1)).
		SetLimit(int64(*params.Limit)).
		SetBatchSize(100)
	
	if len(fields) > 0 {
		opts.SetProjection(utils.BuildMongoProjection(fields))
	}

	cursor, err := r.dealerClientCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoDealerClients []mongoModels.DealerClient
	if err := cursor.All(ctx, &mongoDealerClients); err != nil {
		return nil, err
	}
	return converters.ToDomainDealerClientSlice(mongoDealerClients), nil
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

func (r *MongoDealerClientRepository) Update(ctx context.Context, id string, updates models.DealerClientUpdate) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": updates}
	_, err = r.dealerClientCollection.UpdateByID(ctx, objectID, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoDealerClientRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.dealerClientCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoDealerClientRepository) CheckPhoneExistsForDealer(ctx context.Context, dealerID string, phone string) (bool, error) {
	dealerObjectID, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		return false, err
	}

	

	filter := bson.M{
		"dealer_id":   dealerObjectID,
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


func (r *MongoDealerClientRepository) CreateDealerClientPropertyInterest(ctx context.Context, dealerClientID string, dealerClientPropertyInterest models.DealerClientPropertyInterest) error {
	objectID, err := primitive.ObjectIDFromHex(dealerClientID)
	if err != nil {
		return err
	}
	mongoDealerClientPropertyInterest := converters.ToMongoDealerClientPropertyInterest(dealerClientPropertyInterest)

	_, err = r.dealerClientCollection.UpdateByID(ctx, objectID, bson.M{"$push": bson.M{"properties": mongoDealerClientPropertyInterest}})
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoDealerClientRepository) CheckPropertyInterestExists(ctx context.Context, dealerClientID string, propertyID string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(dealerClientID)
	if err != nil {
		return false, err
	}
	
	filter := bson.M{
		"_id": objectID,
		"properties.property_id": propertyID,
	}

	count, err := r.dealerClientCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}