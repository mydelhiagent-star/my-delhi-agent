package converters

import (
	"myapp/models"
	mongoModels "myapp/mongo_models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToDomainDealerClient(mongoDealerClient mongoModels.DealerClient) models.DealerClient {
	return models.DealerClient{
		ID:                mongoDealerClient.ID.Hex(),
		DealerID:          mongoDealerClient.DealerID.Hex(),
		Name:              mongoDealerClient.Name,
		Phone:             mongoDealerClient.Phone,
		Note:              mongoDealerClient.Note,
		Docs:              ToDomainDealerClientDocs(mongoDealerClient.Docs),
		PropertyInterests: ToDomainDealerClientPropertyInterestSlice(mongoDealerClient.PropertyInterests),
		CreatedAt:         mongoDealerClient.CreatedAt,
		UpdatedAt:         mongoDealerClient.UpdatedAt,
	}
}
func ToDomainDealerClientDocs(mongoDocs []mongoModels.Document) []models.Document {
	docs := make([]models.Document, len(mongoDocs))
	for i, doc := range mongoDocs {
		docs[i] = models.Document{
			URL:  doc.URL,
			Type: doc.Type,
			Size: doc.Size,
		}
	}
	return docs
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
		ID:             mongoDealerClientPropertyInterest.ID.Hex(),
		PropertyID:     mongoDealerClientPropertyInterest.PropertyID.Hex(),
		PropertyNumber: mongoDealerClientPropertyInterest.PropertyNumber,
		Note:           mongoDealerClientPropertyInterest.Note,
		Status:         mongoDealerClientPropertyInterest.Status,
		CreatedAt:      mongoDealerClientPropertyInterest.CreatedAt,
		UpdatedAt:      mongoDealerClientPropertyInterest.UpdatedAt,
	}
}

func ToMongoDealerClientPropertyInterest(dealerClientPropertyInterest models.DealerClientPropertyInterest) mongoModels.DealerClientPropertyInterest {

	propertyObjectID, _ := primitive.ObjectIDFromHex(dealerClientPropertyInterest.PropertyID)
	createdAt, updatedAt := CreationTimestamps()
	return mongoModels.DealerClientPropertyInterest{
		ID:             primitive.NewObjectID(),
		PropertyID:     propertyObjectID,
		PropertyNumber: dealerClientPropertyInterest.PropertyNumber,
		Note:           dealerClientPropertyInterest.Note,
		Status:         dealerClientPropertyInterest.Status,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func ToMongoDealerClientPropertyInterestSlice(dealerClientPropertyInterests []models.DealerClientPropertyInterest) []mongoModels.DealerClientPropertyInterest {
	mongoPropertyInterests := make([]mongoModels.DealerClientPropertyInterest, len(dealerClientPropertyInterests))
	for i, dealerClientPropertyInterest := range dealerClientPropertyInterests {
		mongoPropertyInterests[i] = ToMongoDealerClientPropertyInterest(dealerClientPropertyInterest)
	}
	return mongoPropertyInterests
}

func ToMongoDealerClientUpdate(update models.DealerClientUpdate) mongoModels.DealerClientUpdate {
	return mongoModels.DealerClientUpdate{
		Name:  update.Name,
		Phone: update.Phone,
		Note:  update.Note,
		Docs:  convertDomainDocsToMongoDocs(update.Docs),
	}
}
func convertDomainDocsToMongoDocs(domainDocs *[]models.Document) *[]mongoModels.DocumentUpdate {
	if domainDocs == nil {
		return nil
	}
	mongoDocs := make([]mongoModels.DocumentUpdate, len(*domainDocs))
	for i, doc := range *domainDocs {
		mongoDocs[i] = mongoModels.DocumentUpdate{
			URL:  &doc.URL,
			Type: &doc.Type,
			Size: &doc.Size,
		}
	}
	return &mongoDocs
}

func ToMongoDealerClientDocs(domainDocs []models.Document) []mongoModels.Document {
	mongoDocs := make([]mongoModels.Document, len(domainDocs))
	for i, doc := range domainDocs {
		mongoDocs[i] = mongoModels.Document{
			URL:  doc.URL,
			Type: doc.Type,
			Size: doc.Size,
		}
	}
	return mongoDocs
}
