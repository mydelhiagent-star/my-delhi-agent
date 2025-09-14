package validate

import (
	"errors"
	"myapp/models"
)

func ValidateProperty(property models.Property) error {
	if property.Title == ""{
		return errors.New("title is required")
	}
	if len(property.Title) > 100 {
		return errors.New("title is too long max 100 words")
	}
	if property.Description == ""{
		return errors.New("description is required")
	}
	if len(property.Description) > 500 {
		return errors.New("description is too long max 500 words")
	}
	if property.Address == ""{
		return errors.New("address is required")
	}
	if len(property.Address) > 100 {
		return errors.New("address is too long max 100 words")
	}
	if property.NearestLandmark == ""{
		return errors.New("nearest landmark is required")
	}
	if len(property.NearestLandmark) > 50 {
		return errors.New("nearest landmark is too long max 50 words")
	}
	if property.OwnerName == ""{
		return errors.New("owner name is required")
	}
	if len(property.OwnerName) > 50 {
		return errors.New("owner name is too long max 50 words")
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


func ValidatePropertyUpdate(property models.PropertyUpdate) error {
	if property.Title != nil && *property.Title == "" {
		return errors.New("title is required")
	}
	if property.Title != nil && len(*property.Title) > 100 {
		return errors.New("title is too long max 100 words")
	}
	if property.Description != nil && *property.Description == "" {
		return errors.New("description is required")
	}
	if property.Description != nil && len(*property.Description) > 500 {
		return errors.New("description is too long max 500 words")
	}
	if property.Address != nil && *property.Address == "" {
		return errors.New("address is required")
	}
	if property.Address != nil && len(*property.Address) > 100 {
		return errors.New("address is too long max 100 words")
	}
	if property.NearestLandmark != nil && *property.NearestLandmark == "" {
		return errors.New("nearest landmark is required")
	}
	if property.NearestLandmark != nil && len(*property.NearestLandmark) > 50 {
		return errors.New("nearest landmark is too long max 50 words")
	}
	if property.OwnerName != nil && *property.OwnerName == "" {
		return errors.New("owner name is required")
	}
	if property.OwnerName != nil && len(*property.OwnerName) > 50 {
		return errors.New("owner name is too long max 50 words")
	}
	if err := ValidatePhone(*property.OwnerPhone); err != nil {
		return errors.New("invalid owner phone")
	}
	if property.MinPrice != nil && *property.MinPrice <= 0 {
		return errors.New("min price is required")
	}
	if property.MaxPrice != nil && *property.MaxPrice <= 0 {
		return errors.New("max price is required")
	}
	if property.MinPrice != nil && property.MaxPrice != nil && *property.MinPrice > *property.MaxPrice {
		return errors.New("min price must be less than max price")
	}
	
	if property.PropertyType != nil && *property.PropertyType == "" {
		return errors.New("property type is required")
	}
	if property.Area != nil && *property.Area <= 0 {
		return errors.New("area is required")
	}
	if property.Bedrooms != nil && *property.Bedrooms <= 0 {
		return errors.New("bedrooms is required")
	}
	if property.Bathrooms != nil && *property.Bathrooms <= 0 {
		return errors.New("bathrooms is required")
	}

	return nil
}
