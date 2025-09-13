package mongo_repositories

import (
	"context"
	"myapp/constants"
	"myapp/converters"
	"myapp/models"
	mongoModels "myapp/mongo_models"
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

func (r *MongoDealerRepository) Create(ctx context.Context, dealer models.Dealer) (string, error) {
	mongoDealer, err := converters.ToMongoDealer(dealer)
	if err != nil {
		return "", err
	}

	result, err := r.dealerCollection.InsertOne(ctx, mongoDealer)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoDealerRepository) GetByID(ctx context.Context, id string) (models.Dealer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Dealer{}, err
	}

	var mongoDealer mongoModels.Dealer
	err = r.dealerCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&mongoDealer)
	if err != nil {
		return models.Dealer{}, err
	}

	return converters.ToDomainDealer(mongoDealer), nil
}

func (r *MongoDealerRepository) GetByPhone(ctx context.Context, phone string) (models.Dealer, error) {
	var mongoDealer mongoModels.Dealer
	err := r.dealerCollection.FindOne(ctx, bson.M{"phone": phone}).Decode(&mongoDealer)
	if err != nil {
		return models.Dealer{}, err
	}

	return converters.ToDomainDealer(mongoDealer), nil
}

func (r *MongoDealerRepository) GetByEmail(ctx context.Context, email string) (models.Dealer, error) {
	var mongoDealer mongoModels.Dealer
	err := r.dealerCollection.FindOne(ctx, bson.M{"email": email}).Decode(&mongoDealer)
	if err != nil {
		return models.Dealer{}, err
	}

	return converters.ToDomainDealer(mongoDealer), nil
}

func (r *MongoDealerRepository) GetAll(ctx context.Context) ([]models.Dealer, error) {
	cursor, err := r.dealerCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoDealers []mongoModels.Dealer
	if err := cursor.All(ctx, &mongoDealers); err != nil {
		return nil, err
	}

	return converters.ToDomainDealerSlice(mongoDealers), nil
}

func (r *MongoDealerRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": updates}
	_, err = r.dealerCollection.UpdateByID(ctx, objectID, update)
	return err
}

func (r *MongoDealerRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.dealerCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *MongoDealerRepository) Exists(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	count, err := r.dealerCollection.CountDocuments(ctx, bson.M{"_id": objectID})
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

	var mongoDealers []mongoModels.Dealer
	if err := cursor.All(ctx, &mongoDealers); err != nil {
		return nil, err
	}

	return converters.ToDomainDealerSlice(mongoDealers), nil
}

func (r *MongoDealerRepository) GetLocationsWithSubLocations(ctx context.Context) ([]models.LocationWithSubLocations, error) {
	cursor, err := r.dealerCollection.Find(
		ctx,
		bson.M{"location": bson.M{"$in": constants.Locations}},
		options.Find().SetProjection(bson.M{
			"location":     1,
			"sub_location": 1,
			"_id":          0,
		}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Map: location -> set of sublocations
	locationMap := make(map[string]map[string]struct{})

	for cursor.Next(ctx) {
		var mongoDealer mongoModels.Dealer
		if err := cursor.Decode(&mongoDealer); err != nil {
			return nil, err
		}

		if mongoDealer.Location != "" && mongoDealer.SubLocation != "" {
			if _, ok := locationMap[mongoDealer.Location]; !ok {
				locationMap[mongoDealer.Location] = make(map[string]struct{})
			}
			locationMap[mongoDealer.Location][mongoDealer.SubLocation] = struct{}{}
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