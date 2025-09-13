package converters

import (
	"myapp/models"
	mongoModels "myapp/mongo_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func ToMongoDealer(dealer models.Dealer) (mongoModels.Dealer, error) {
	mongoDealer := mongoModels.Dealer{
		Name:          dealer.Name,
		Phone:         dealer.Phone,
		Password:      dealer.Password,
		Email:         dealer.Email,
		OfficeAddress: dealer.OfficeAddress,
		ShopName:      dealer.ShopName,
		Location:      dealer.Location,
		SubLocation:   dealer.SubLocation,
		CreatedAt:     dealer.CreatedAt,
		UpdatedAt:     dealer.UpdatedAt,
	}

	if dealer.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(dealer.ID)
		if err != nil {
			return mongoModels.Dealer{}, err
		}
		mongoDealer.ID = objectID
	}

	return mongoDealer, nil
}

func ToDomainDealer(mongoDealer mongoModels.Dealer) models.Dealer {
	return models.Dealer{
		ID:            mongoDealer.ID.Hex(),
		Name:          mongoDealer.Name,
		Phone:         mongoDealer.Phone,
		Password:      mongoDealer.Password,
		Email:         mongoDealer.Email,
		OfficeAddress: mongoDealer.OfficeAddress,
		ShopName:      mongoDealer.ShopName,
		Location:      mongoDealer.Location,
		SubLocation:   mongoDealer.SubLocation,
		CreatedAt:     mongoDealer.CreatedAt,
		UpdatedAt:     mongoDealer.UpdatedAt,
	}
}



func ToDomainDealerSlice(mongoDealers []mongoModels.Dealer) []models.Dealer {
	dealers := make([]models.Dealer, len(mongoDealers))
	for i, mongoDealer := range mongoDealers {
		dealers[i] = ToDomainDealer(mongoDealer)
	}
	return dealers
}