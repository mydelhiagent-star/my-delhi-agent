package converters

import (
	"myapp/models"
	mongoModels "myapp/mongo_models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



// Convert domain model to MongoDB model
func ToMongoProperty(property models.Property) (mongoModels.Property, error) {
	mongoProperty := mongoModels.Property{
		PropertyNumber:  property.PropertyNumber,
		Title:           property.Title,
		Description:     property.Description,
		Address:         property.Address,
		MinPrice:        property.MinPrice,
		MaxPrice:        property.MaxPrice,
		Photos:          property.Photos,
		Videos:          property.Videos,
		OwnerName:       property.OwnerName,
		OwnerPhone:      property.OwnerPhone,
		NearestLandmark: property.NearestLandmark,
		IsDeleted:       property.IsDeleted,
		Sold:            property.Sold,
		SoldPrice:       property.SoldPrice,
		SoldDate:        property.SoldDate,
		Area:            property.Area,
		Bedrooms:        property.Bedrooms,
		Bathrooms:       property.Bathrooms,
		PropertyType:    property.PropertyType,
		CreatedAt:       property.CreatedAt,
		UpdatedAt:       property.UpdatedAt,
	}

	// Convert string ID to ObjectID
	if property.ID != "" {
		objectID, err := primitive.ObjectIDFromHex(property.ID)
		if err != nil {
			return mongoModels.Property{}, err
		}
		mongoProperty.ID = objectID
	}

	if property.DealerID != "" {
		dealerObjectID, err := primitive.ObjectIDFromHex(property.DealerID)
		if err != nil {
			return mongoModels.Property{}, err
		}
		mongoProperty.DealerID = dealerObjectID
	}

	return mongoProperty, nil
}


func ToDomainProperty(mongoProperty mongoModels.Property) models.Property {
	return models.Property{
		ID:              mongoProperty.ID.Hex(),
		PropertyNumber:  mongoProperty.PropertyNumber,
		DealerID:        mongoProperty.DealerID.Hex(),
		Title:           mongoProperty.Title,
		Description:     mongoProperty.Description,
		Address:         mongoProperty.Address,
		MinPrice:        mongoProperty.MinPrice,
		MaxPrice:        mongoProperty.MaxPrice,
		Photos:          mongoProperty.Photos,
		Videos:          mongoProperty.Videos,
		OwnerName:       mongoProperty.OwnerName,
		OwnerPhone:      mongoProperty.OwnerPhone,
		NearestLandmark: mongoProperty.NearestLandmark,
		IsDeleted:       mongoProperty.IsDeleted,
		Sold:            mongoProperty.Sold,
		SoldPrice:       mongoProperty.SoldPrice,
		SoldDate:        mongoProperty.SoldDate,
		Area:            mongoProperty.Area,
		Bedrooms:        mongoProperty.Bedrooms,
		Bathrooms:       mongoProperty.Bathrooms,
		PropertyType:    mongoProperty.PropertyType,
		CreatedAt:       mongoProperty.CreatedAt,
		UpdatedAt:       mongoProperty.UpdatedAt,
	}
}

// Convert slice of MongoDB models to domain models
func ToDomainPropertySlice(mongoProperties []mongoModels.Property) []models.Property {
	properties := make([]models.Property, len(mongoProperties))
	for i, mongoProperty := range mongoProperties {
		properties[i] = ToDomainProperty(mongoProperty)
	}
	return properties
}

// converters/property.go
func ToMongoPropertyUpdate(update models.PropertyUpdate) mongoModels.PropertyUpdate {
    return mongoModels.PropertyUpdate{
        Title:           update.Title,
        Address:         update.Address,
        NearestLandmark: update.NearestLandmark,
        SoldBy:          update.SoldBy,
        MinPrice:        update.MinPrice,
        MaxPrice:        update.MaxPrice,
        Description:     update.Description,
        Photos:          update.Photos,
        Videos:          update.Videos,
        OwnerName:       update.OwnerName,
        OwnerPhone:      update.OwnerPhone,
        Sold:            update.Sold,
        IsDeleted:       update.IsDeleted,
        SoldPrice:       update.SoldPrice,
        SoldDate:        update.SoldDate,
        Area:            update.Area,
        Bedrooms:        update.Bedrooms,
        Bathrooms:       update.Bathrooms,
        PropertyType:    update.PropertyType,
        UpdatedAt:       update.UpdatedAt,
    }
}
