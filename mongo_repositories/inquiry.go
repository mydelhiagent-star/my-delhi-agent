package mongo_repositories

import (
	"context"
	"time"

	"myapp/converters"
	"myapp/models"
	mongoModels "myapp/mongo_models"
	"myapp/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInquiryRepository struct {
	inquiryCollection *mongo.Collection
}

func NewMongoInquiryRepository(inquiryCollection *mongo.Collection) *MongoInquiryRepository {
	return &MongoInquiryRepository{
		inquiryCollection: inquiryCollection,
	}
}

func (r *MongoInquiryRepository) Create(ctx context.Context, inquiry models.Inquiry) (models.Inquiry, error) {
	mongoInquiry := converters.ToMongoInquiry(inquiry)
	mongoInquiry.CreatedAt = time.Now()
	mongoInquiry.UpdatedAt = time.Now()

	result, err := r.inquiryCollection.InsertOne(ctx, mongoInquiry)
	if err != nil {
		return models.Inquiry{}, err
	}

	mongoInquiry.ID = result.InsertedID.(primitive.ObjectID)
	return converters.ToDomainInquiry(mongoInquiry), nil
}

func (r *MongoInquiryRepository) GetByID(ctx context.Context, id string) (models.Inquiry, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Inquiry{}, err
	}

	var mongoInquiry mongoModels.Inquiry
	err = r.inquiryCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&mongoInquiry)
	if err != nil {
		return models.Inquiry{}, err
	}

	return converters.ToDomainInquiry(mongoInquiry), nil
}

func (r *MongoInquiryRepository) GetAll(ctx context.Context, params models.InquiryQueryParams) ([]models.Inquiry, error) {
	params.SetDefaults()

	filter := utils.BuildMongoFilter(params)

	var mongoInquiries []mongoModels.Inquiry
	var cursor *mongo.Cursor
	var err error

	if *params.Aggregation {
		pipeline := utils.BuildAggregationPipeline(filter, *params.Sort, getSortOrder(*params.Order), int64((*params.Page-1)*(*params.Limit)), int64(*params.Limit), bson.M{})
		cursor, err = r.inquiryCollection.Aggregate(ctx, pipeline)
	} else {
		opts := options.Find().
			SetSort(bson.D{{Key: *params.Sort, Value: getSortOrder(*params.Order)}}).
			SetSkip(int64((*params.Page - 1) * (*params.Limit))).
			SetLimit(int64(*params.Limit))
		cursor, err = r.inquiryCollection.Find(ctx, filter, opts)
	}

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &mongoInquiries)
	if err != nil {
		return nil, err
	}

	return converters.ToDomainInquirySlice(mongoInquiries), nil
}

func (r *MongoInquiryRepository) Update(ctx context.Context, id string, updates models.InquiryUpdate) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	mongoUpdate := converters.ToMongoInquiryUpdate(updates)
	updateDoc := utils.BuildUpdateDocument(mongoUpdate)

	_, err = r.inquiryCollection.UpdateByID(ctx, objectID, bson.M{"$set": updateDoc})
	return err
}

func (r *MongoInquiryRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.inquiryCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func getSortOrder(order string) int {
	if order == "asc" {
		return 1
	}
	return -1
}
