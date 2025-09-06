package validate

import (
	"errors"
	"myapp/models"
)

func ValidateProperty(property models.Property) error {
	if property.Title == "" || len(property.Title) < 2 {
		return errors.New("title is required")
	}
	if property.Description == "" || len(property.Description) < 2 {
		return errors.New("description is required")
	}
	if property.Address == "" || len(property.Address) < 2 {
		return errors.New("address is required")
	}
	if property.NearestLandmark == "" || len(property.NearestLandmark) < 2 {
		return errors.New("nearest landmark is required")
	}
	if property.OwnerName == "" || len(property.OwnerName) < 2 {
		return errors.New("owner name is required")
	}
	if property.Area <= 0 {
		return errors.New("area is required")
	}
	if property.Bedrooms <= 0 {
		return errors.New("bedrooms is required")
	}
	if property.Bathrooms <= 0 {
		return errors.New("bathrooms is required")
	}
	if property.PropertyType == "" {
		return errors.New("property type is required")
	}
	if property.MinPrice <= 0 {
		return errors.New("min price is required")
	}
	if property.MaxPrice <= 0 {
		return errors.New("max price is required")
	}
	if property.MinPrice > property.MaxPrice {
		return errors.New("min price must be less than max price")
	}

	if err := ValidatePhone(property.OwnerPhone); err != nil {
		return errors.New("invalid owner phone")
	}

	return nil
}

