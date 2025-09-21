package converters

import (
	"myapp/models"
	mongoModels "myapp/mongo_models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToDomainInquiry(mongoInquiry mongoModels.Inquiry) models.Inquiry {
	inquiry := models.Inquiry{
		ID:          mongoInquiry.ID.Hex(),
		Source:      mongoInquiry.Source,
		Name:        mongoInquiry.Name,
		Phone:       mongoInquiry.Phone,
		Requirement: mongoInquiry.Requirement,
		CreatedAt:   mongoInquiry.CreatedAt,
		UpdatedAt:   mongoInquiry.UpdatedAt,
	}

	if mongoInquiry.DealerID != nil {
		dealerID := mongoInquiry.DealerID.Hex()
		inquiry.DealerID = &dealerID
	}

	return inquiry
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

	var dealerObjectID *primitive.ObjectID
	if inquiry.DealerID != nil && *inquiry.DealerID != "" {
		if id, err := primitive.ObjectIDFromHex(*inquiry.DealerID); err == nil {
			dealerObjectID = &id
		}
	}

	return mongoModels.Inquiry{
		ID:          objectID,
		DealerID:    dealerObjectID,
		Source:      inquiry.Source,
		Name:        inquiry.Name,
		Phone:       inquiry.Phone,
		Requirement: inquiry.Requirement,
		CreatedAt:   inquiry.CreatedAt,
		UpdatedAt:   inquiry.UpdatedAt,
	}
}

func ToMongoInquiryUpdate(update models.InquiryUpdate) mongoModels.InquiryUpdate {
	var dealerObjectID *primitive.ObjectID
	if update.DealerID != nil && *update.DealerID != "" {
		if id, err := primitive.ObjectIDFromHex(*update.DealerID); err == nil {
			dealerObjectID = &id
		}
	}

	return mongoModels.InquiryUpdate{
		DealerID:    dealerObjectID,
		Source:      update.Source,
		Name:        update.Name,
		Phone:       update.Phone,
		Requirement: update.Requirement,
		UpdatedAt:   update.UpdatedAt,
	}
}
