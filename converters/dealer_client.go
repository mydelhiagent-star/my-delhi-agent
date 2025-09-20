package converters

import (
	"myapp/models"
	mongoModels "myapp/mongo_models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToDomainDealerClient(mongoDealerClient mongoModels.DealerClient) models.DealerClient {
	return models.DealerClient{
		ID: mongoDealerClient.ID.Hex(),
		DealerID: mongoDealerClient.DealerID.Hex(),
		Name: mongoDealerClient.Name,
		Phone: mongoDealerClient.Phone,
		Note: mongoDealerClient.Note,
		PropertyInterests: ToDomainDealerClientPropertyInterestSlice(mongoDealerClient.PropertyInterests),
		CreatedAt: mongoDealerClient.CreatedAt,
		UpdatedAt: mongoDealerClient.UpdatedAt,
	}
}

func ToDomainDealerClientSlice(mongoDealerClients []mongoModels.DealerClient) []models.DealerClient {
	dealerClients := make([]models.DealerClient, len(mongoDealerClients))
	for i, mongoDealerClient := range mongoDealerClients {
		dealerClients[i] = ToDomainDealerClient(mongoDealerClient)
	}
	return dealerClients
}

func ToDomainDealerClientPropertyInterestSlice(mongoDealerClientPropertyInterests []mongoModels.DealerClientPropertyInterest) []models.DealerClientPropertyInterest {	
	dealerClientPropertyInterests := make([]models.DealerClientPropertyInterest, len(mongoDealerClientPropertyInterests))
	for i, mongoDealerClientPropertyInterest := range mongoDealerClientPropertyInterests {
		dealerClientPropertyInterests[i] = ToDomainDealerClientPropertyInterest(mongoDealerClientPropertyInterest)
	}
	return dealerClientPropertyInterests
}

func ToDomainDealerClientPropertyInterest(mongoDealerClientPropertyInterest mongoModels.DealerClientPropertyInterest) models.DealerClientPropertyInterest {
	return models.DealerClientPropertyInterest{
		ID: mongoDealerClientPropertyInterest.ID.Hex(),
		PropertyID: mongoDealerClientPropertyInterest.PropertyID.Hex(),
		Note: mongoDealerClientPropertyInterest.Note,
		Status: mongoDealerClientPropertyInterest.Status,
		CreatedAt: mongoDealerClientPropertyInterest.CreatedAt,
		UpdatedAt: mongoDealerClientPropertyInterest.UpdatedAt,
	}
}


func ToMongoDealerClientPropertyInterest(dealerClientPropertyInterest models.DealerClientPropertyInterest) mongoModels.DealerClientPropertyInterest {
    
    propertyObjectID, _ := primitive.ObjectIDFromHex(dealerClientPropertyInterest.PropertyID)
    createdAt, updatedAt := CreationTimestamps()
    return mongoModels.DealerClientPropertyInterest{
        ID:         primitive.NewObjectID(), 
        PropertyID: propertyObjectID,       
        Note:       dealerClientPropertyInterest.Note,
        Status:     dealerClientPropertyInterest.Status,
        CreatedAt:  createdAt,              
        UpdatedAt:  updatedAt,             
    }
}

func ToMongoDealerClientPropertyInterestSlice(dealerClientPropertyInterests []models.DealerClientPropertyInterest) []mongoModels.DealerClientPropertyInterest {
	mongoPropertyInterests := make([]mongoModels.DealerClientPropertyInterest, len(dealerClientPropertyInterests))
	for i, dealerClientPropertyInterest := range dealerClientPropertyInterests {
		mongoPropertyInterests[i] = ToMongoDealerClientPropertyInterest(dealerClientPropertyInterest)
	}
	return mongoPropertyInterests
}