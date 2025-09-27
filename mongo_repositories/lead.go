package mongo_repositories

import (
	"context"
	"errors"
	"fmt"
	"myapp/converters"
	"myapp/models"
	mongoModels "myapp/mongo_models"
	"myapp/repositories"
	"myapp/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoLeadRepository struct {
	leadCollection     *mongo.Collection
	propertyCollection *mongo.Collection
}

func NewMongoLeadRepository(leadCollection, propertyCollection *mongo.Collection) repositories.LeadRepository {
	return &MongoLeadRepository{
		leadCollection:     leadCollection,
		propertyCollection: propertyCollection,
	}
}

func (r *MongoLeadRepository) Create(ctx context.Context, lead models.Lead) (string, error) {
	// Convert models.LeadData to mongoModels.Lead
	mongoLead := mongoModels.Lead{
		Name:         lead.Name,
		Phone:        lead.Phone,
		Requirement:  lead.Requirement,
		AadharNumber: lead.AadharNumber,
		AadharPhoto:  lead.AadharPhoto,
	}

	res, err := r.leadCollection.InsertOne(ctx, mongoLead)
	if err != nil {
		return "", err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to convert inserted ID")
	}
	return id.Hex(), nil
}

func (r *MongoLeadRepository) GetByID(ctx context.Context, id string) (models.Lead, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Lead{}, err
	}

	var mongoLead mongoModels.Lead
	err = r.leadCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&mongoLead)
	if err != nil {
		return models.Lead{}, err
	}

	// Convert mongoModels.Lead to models.LeadData
	return models.Lead{
		ID:           mongoLead.ID.Hex(),
		Name:         mongoLead.Name,
		Phone:        mongoLead.Phone,
		Requirement:  mongoLead.Requirement,
		AadharNumber: mongoLead.AadharNumber,
		AadharPhoto:  mongoLead.AadharPhoto,
	}, nil
}

func (r *MongoLeadRepository) GetAll(ctx context.Context) ([]models.Lead, error) {
	cursor, err := r.leadCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var mongoLeads []mongoModels.Lead
	if err = cursor.All(ctx, &mongoLeads); err != nil {
		return nil, err
	}

	// Convert mongoModels.Lead to models.LeadData
	var leads []models.Lead
	for _, mongoLead := range mongoLeads {
		leads = append(leads, models.Lead{
			ID:           mongoLead.ID.Hex(),
			Name:         mongoLead.Name,
			Phone:        mongoLead.Phone,
			Requirement:  mongoLead.Requirement,
			AadharNumber: mongoLead.AadharNumber,
			AadharPhoto:  mongoLead.AadharPhoto,
		})
	}
	return leads, nil
}

func (r *MongoLeadRepository) GetByDealerID(ctx context.Context, dealerID string) ([]models.Lead, error) {
	dealerObjectID, err := primitive.ObjectIDFromHex(dealerID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.leadCollection.Find(ctx, bson.M{"dealer_id": dealerObjectID})
	if err != nil {
		return nil, err
	}

	var mongoLeads []mongoModels.Lead
	if err = cursor.All(ctx, &mongoLeads); err != nil {
		return nil, err
	}

	// Convert mongoModels.Lead to models.LeadData
	var leads []models.Lead
	for _, mongoLead := range mongoLeads {
		leads = append(leads, models.Lead{
			ID:           mongoLead.ID.Hex(),
			Name:         mongoLead.Name,
			Phone:        mongoLead.Phone,
			Requirement:  mongoLead.Requirement,
			AadharNumber: mongoLead.AadharNumber,
			AadharPhoto:  mongoLead.AadharPhoto,
		})
	}
	return leads, nil
}

func (r *MongoLeadRepository) GetLeads(ctx context.Context, params models.LeadQueryParams) ([]models.Lead, error) {
	params.SetDefaults()
	
	filter := utils.BuildMongoFilter(params)
	if *params.Aggregation {
		return r.getLeadsByAggregation(ctx, filter, params, []string{})
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
		SetLimit(int64(limit))
	
	cursor, err := r.leadCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoLeads []mongoModels.Lead
	if err := cursor.All(ctx, &mongoLeads); err != nil {
		return nil, err
	}

	return converters.ToDomainLeadSlice(mongoLeads), nil
}

func (r *MongoLeadRepository) getLeadsByAggregation(ctx context.Context,filter bson.M, params models.LeadQueryParams,fields []string) ([]models.Lead, error) {
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
	cursor, err := r.leadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var mongoLeads []mongoModels.Lead
	if err := cursor.All(ctx, &mongoLeads); err != nil {
		return nil, err
	}
	return converters.ToDomainLeadSlice(mongoLeads), nil
}

	


func (r *MongoLeadRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": updates}
	_, err = r.leadCollection.UpdateByID(ctx, objectID, update)
	return err
}

func (r *MongoLeadRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.leadCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoLeadRepository) AddPropertyInterest(ctx context.Context, leadID string, propertyInterest models.PropertyInterest) error {
	leadObjectID, err := primitive.ObjectIDFromHex(leadID)
	if err != nil {
		return err
	}

	propertyObjectID, err := primitive.ObjectIDFromHex(propertyInterest.PropertyID)
	if err != nil {
		return err
	}

	dealerObjectID, err := primitive.ObjectIDFromHex(propertyInterest.DealerID)
	if err != nil {
		return err
	}

	// Convert models.PropertyInterestData to mongoModels.PropertyInterest
	mongoPropertyInterest := mongoModels.PropertyInterest{
		PropertyID: propertyObjectID,
		PropertyNumber: propertyInterest.PropertyNumber,
		DealerID:   dealerObjectID,
		Status:     propertyInterest.Status,
		Note:       propertyInterest.Note,
	}

	// Check if property already exists for this lead
	filter := bson.M{
		"_id":                    leadObjectID,
		"properties.property_id": propertyObjectID,
	}

	var existingLead mongoModels.Lead
	err = r.leadCollection.FindOne(ctx, filter).Decode(&existingLead)
	if err == nil {
		return errors.New("property already added to this lead")
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("database error checking property: %w", err)
	}

	// Add property interest
	update := bson.M{
		"$push": bson.M{
			"properties": mongoPropertyInterest,
		},
	}

	_, err = r.leadCollection.UpdateOne(ctx, bson.M{"_id": leadObjectID}, update)
	return err
}

func (r *MongoLeadRepository) UpdatePropertyInterest(ctx context.Context, leadID, propertyID string, status, note string) error {
	leadObjectID, err := primitive.ObjectIDFromHex(leadID)
	if err != nil {
		return err
	}

	propertyObjectID, err := primitive.ObjectIDFromHex(propertyID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":                    leadObjectID,
		"properties.property_id": propertyObjectID,
	}

	update := bson.M{
		"$set": bson.M{
			"properties.$.status": status,
			"properties.$.note":   note,
		},
	}

	result, err := r.leadCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}






func (r *MongoLeadRepository) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	count, err := r.leadCollection.CountDocuments(ctx, bson.M{"phone": phone})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
