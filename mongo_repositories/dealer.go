package mongo_repository

import (
	"context"
	"myapp/constants"
	"myapp/models"
	"myapp/repositories"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDealerRepository struct {
	dealerCollection *mongo.Collection
}

func NewMongoDealerRepository(dealerCollection *mongo.Collection) repositories.DealerRepository {
	return &MongoDealerRepository{
		dealerCollection: dealerCollection,
	}
}

func (r *MongoDealerRepository) Create(ctx context.Context, dealer models.Dealer) (primitive.ObjectID, error) {
	result, err := r.dealerCollection.InsertOne(ctx, dealer)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *MongoDealerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Dealer, error) {
	var dealer models.Dealer
	err := r.dealerCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&dealer)
	if err != nil {
		return nil, err
	}
	return &dealer, nil
}

func (r *MongoDealerRepository) GetByPhone(ctx context.Context, phone string) (*models.Dealer, error) {
	var dealer models.Dealer
	err := r.dealerCollection.FindOne(ctx, bson.M{"phone": phone}).Decode(&dealer)
	if err != nil {
		return nil, err
	}
	return &dealer, nil
}

func (r *MongoDealerRepository) GetByEmail(ctx context.Context, email string) (*models.Dealer, error) {
	var dealer models.Dealer
	err := r.dealerCollection.FindOne(ctx, bson.M{"email": email}).Decode(&dealer)
	if err != nil {
		return nil, err
	}
	return &dealer, nil
}

func (r *MongoDealerRepository) GetAll(ctx context.Context) ([]models.Dealer, error) {
	cursor, err := r.dealerCollection.Find(ctx, bson.M{})
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

func (r *MongoDealerRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	update := bson.M{"$set": updates}
	_, err := r.dealerCollection.UpdateByID(ctx, id, update)
	return err
}

func (r *MongoDealerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.dealerCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoDealerRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	count, err := r.dealerCollection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MongoDealerRepository) GetByLocation(ctx context.Context, subLocation string) ([]models.Dealer, error) {
	filter := bson.M{
		"sub_location": subLocation,
	}

	cursor, err := r.dealerCollection.Find(ctx, filter)
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

func (r *MongoDealerRepository) GetLocationsWithSubLocations(ctx context.Context) ([]models.LocationWithSubLocations, error) {
	cursor, err := r.dealerCollection.Find(
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
		var dealer models.Dealer
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
