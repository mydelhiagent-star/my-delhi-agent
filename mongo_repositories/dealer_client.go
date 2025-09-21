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
		DealerID:          dealerObjectID,
		Name:              dealerClient.Name,
		Phone:             dealerClient.Phone,
		Note:              dealerClient.Note,
		Docs:              dealerClient.Docs,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		PropertyInterests: mongoPropertyInterests,
	}

	result, err := r.dealerClientCollection.InsertOne(ctx, mongoDealerClient)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoDealerClientRepository) GetByID(ctx context.Context, id string) (models.DealerClient, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.DealerClient{}, err
	}
	var dealerClient models.DealerClient
	err = r.dealerClientCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&dealerClient)
	if err != nil {
		return models.DealerClient{}, err
	}
	return dealerClient, nil
}

func (r *MongoDealerClientRepository) GetDealerClients(ctx context.Context, params models.DealerClientQueryParams, fields []string) ([]models.DealerClient, error) {

	filter := utils.BuildMongoFilter(params)

	if *params.Aggregation {
		return r.getDealerClientsWithAggregation(ctx, filter, params, fields)
	}

	skip := (*params.Page - 1) * *params.Limit
	limit := *params.Limit
	sortValue := 1
	if *params.Order == "desc" {
		sortValue = -1
	}

	opts := options.Find().
		SetSort(bson.M{*params.Sort: sortValue}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit) + 1).
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

func (r *MongoDealerClientRepository) getDealerClientsWithAggregation(ctx context.Context, filter bson.M, params models.DealerClientQueryParams, fields []string) ([]models.DealerClient, error) {
	skip := int64((*params.Page - 1) * *params.Limit)
	limit := int64(*params.Limit + 1)
	sort := *params.Sort
	sortValue := 1
	if *params.Order == "desc" {
		sortValue = -1
	}
	var projection bson.M
	if len(fields) > 0 {
		projection = utils.BuildMongoProjection(fields)
	}

	pipeline := utils.BuildAggregationPipeline(filter, sort, sortValue, skip, limit, projection)

	cursor, err := r.dealerClientCollection.Aggregate(ctx, pipeline)
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

func (r *MongoDealerClientRepository) Update(ctx context.Context, id string, updates models.DealerClientUpdate) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	mongoUpdate := converters.ToMongoDealerClientUpdate(updates)
	updateDoc := utils.BuildUpdateDocument(mongoUpdate)
	updateDoc["updated_at"] = time.Now()
	update := bson.M{"$set": updateDoc}
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
		"dealer_id": dealerObjectID,
		"phone":     phone,
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

	propertyObjectID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"_id":                    objectID,
		"properties.property_id": propertyObjectID,
	}

	count, err := r.dealerClientCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MongoDealerClientRepository) UpdateDealerClientPropertyInterest(ctx context.Context, dealerClientID string, propertyInterestID string, update models.DealerClientPropertyInterestUpdate) error {
	objectID, err := primitive.ObjectIDFromHex(dealerClientID)
	if err != nil {
		return err
	}
	propertyInterestObjectID, err := primitive.ObjectIDFromHex(propertyInterestID)
	if err != nil {
		return err
	}

	// ✅ Correct filter
	filter := bson.M{
		"_id":                    objectID,
		"properties.property_id": propertyInterestObjectID,
	}

	// ✅ Conditional updates
	updateDoc := bson.M{}
	if update.Note != nil {
		updateDoc["properties.$.note"] = *update.Note
	}
	if update.Status != nil {
		updateDoc["properties.$.status"] = *update.Status
	}
	updateDoc["updated_at"] = time.Now()

	if len(updateDoc) > 0 {
		_, err = r.dealerClientCollection.UpdateOne(ctx, filter, bson.M{"$set": updateDoc})
		return err
	}

	return nil
}

func (r *MongoDealerClientRepository) DeleteDealerClientPropertyInterest(ctx context.Context, dealerClientID string, propertyInterestID string) error {
	objectID, err := primitive.ObjectIDFromHex(dealerClientID)
	if err != nil {
		return err
	}
	propertyInterestObjectID, err := primitive.ObjectIDFromHex(propertyInterestID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id":                    objectID,
		"properties.property_id": propertyInterestObjectID,
	}
	_, err = r.dealerClientCollection.UpdateOne(ctx, filter, bson.M{"$pull": bson.M{"properties": bson.M{"property_id": propertyInterestObjectID}}})
	return err
}
