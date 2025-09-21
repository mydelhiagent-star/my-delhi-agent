package converters

import (
	"myapp/models"
	mongoModels "myapp/mongo_models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToDomainInquiry(mongoInquiry mongoModels.Inquiry) models.Inquiry {
	return models.Inquiry{
		ID:          mongoInquiry.ID.Hex(),
		Name:        mongoInquiry.Name,
		Phone:       mongoInquiry.Phone,
		Requirement: mongoInquiry.Requirement,
		CreatedAt:   mongoInquiry.CreatedAt,
		UpdatedAt:   mongoInquiry.UpdatedAt,
	}
}

func ToDomainInquirySlice(mongoInquiries []mongoModels.Inquiry) []models.Inquiry {
	inquiries := make([]models.Inquiry, len(mongoInquiries))
	for i, mongoInquiry := range mongoInquiries {
		inquiries[i] = ToDomainInquiry(mongoInquiry)
	}
	return inquiries
}

func ToMongoInquiry(inquiry models.Inquiry) mongoModels.Inquiry {
	var objectID primitive.ObjectID
	if inquiry.ID != "" {
		objectID, _ = primitive.ObjectIDFromHex(inquiry.ID)
	}

	return mongoModels.Inquiry{
		ID:          objectID,
		Name:        inquiry.Name,
		Phone:       inquiry.Phone,
		Requirement: inquiry.Requirement,
		CreatedAt:   inquiry.CreatedAt,
		UpdatedAt:   inquiry.UpdatedAt,
	}
}

func ToMongoInquiryUpdate(update models.InquiryUpdate) mongoModels.InquiryUpdate {
	return mongoModels.InquiryUpdate{
		Name:        update.Name,
		Phone:       update.Phone,
		Requirement: update.Requirement,
		UpdatedAt:   update.UpdatedAt,
	}
}
