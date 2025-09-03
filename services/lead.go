package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"myapp/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeadService struct {
	LeadCollection *mongo.Collection
}

func (s *LeadService) CreateLead(ctx context.Context, lead models.Lead) (primitive.ObjectID, error) {

	existingLead := models.Lead{}
	err := s.LeadCollection.FindOne(ctx, bson.M{"phone": lead.Phone}).Decode(&existingLead)
	if err == nil {
		return primitive.NilObjectID, errors.New("lead with phone already exists")
	} else if err != mongo.ErrNoDocuments {
		return primitive.NilObjectID, errors.New("database error checking phone")
	}

	res, err := s.LeadCollection.InsertOne(ctx, lead)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, err
	}
	return id, nil
}

func (s *LeadService) GetLeadByID(ctx context.Context, id primitive.ObjectID) (models.Lead, error) {

	var lead models.Lead
	err := s.LeadCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&lead)
	return lead, err
}

func (s *LeadService) GetAllLeads(ctx context.Context) ([]models.Lead, error) {
	cursor, err := s.LeadCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var leads []models.Lead
	if err = cursor.All(ctx, &leads); err != nil {
		return nil, err
	}
	return leads, nil
}

func (s *LeadService) GetAllLeadsByDealerID(ctx context.Context, dealerID primitive.ObjectID) ([]models.Lead, error) {
	cursor, err := s.LeadCollection.Find(ctx, bson.M{"dealer_id": dealerID})
	if err != nil {
		return nil, err
	}

	var leads []models.Lead
	if err = cursor.All(ctx, &leads); err != nil {
		return nil, err
	}
	return leads, nil

}

func (s *LeadService) UpdateLead(ctx context.Context, id primitive.ObjectID, updateData map[string]interface{}) error {

	update := bson.M{
		"$set": updateData,
	}
	_, err := s.LeadCollection.UpdateByID(ctx, id, update)
	return err
}

func (s *LeadService) DeleteLead(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.LeadCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *LeadService) AddPropertyInterest(ctx context.Context, leadID primitive.ObjectID, propertyInterest models.PropertyInterest) error {
	// Set status
	propertyInterest.Status = models.LeadStatusView
	propertyInterest.CreatedAt = time.Now()

	// ← CHECK if property already exists for this lead
	var existingLead models.Lead
	err := s.LeadCollection.FindOne(ctx, bson.M{
		"_id":                    leadID,
		"properties.property_id": propertyInterest.PropertyID,
	}).Decode(&existingLead)

	if err == nil {
		// ← Property already exists for this lead
		return errors.New("property already added to this lead")
	} else if err != mongo.ErrNoDocuments {
		// ← Database error occurred
		return fmt.Errorf("database error checking property: %w", err)
	}

	// ← Property doesn't exist, add it
	_, err = s.LeadCollection.UpdateOne(ctx,
		bson.M{"_id": leadID},
		bson.M{"$push": bson.M{"properties": propertyInterest}},
	)

	return err
}

func (s *LeadService) SearchLeads(ctx context.Context, filter bson.M, page, limit int, fields []string) ([]models.Lead, error) {
	// ← CALCULATE skip value for pagination
	skip := (page - 1) * limit

	// ← BUILD aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match leads based on filter
		{{Key: "$match", Value: filter}},

		// Stage 2: Lookup properties
		{{Key: "$lookup", Value: bson.M{
			"from":         "property",
			"localField":   "properties.property_id",
			"foreignField": "_id",
			"as":           "property_details",
		}}},

		// Stage 3: Filter properties based on deleted + sold rules
		{{Key: "$addFields", Value: bson.M{
			"properties": bson.M{
				"$filter": bson.M{
					"input": "$properties",
					"as":    "prop",
					"cond": bson.M{
						"$and": bson.A{
							// 1. Exclude deleted properties
							bson.M{"$ne": bson.A{"$$prop.is_deleted", true}},

							// 2. Sold logic
							bson.M{
								"$or": bson.A{
									// Case A: If status = converted → always include
									bson.M{"$eq": bson.A{"$status", "converted"}},

									// Case B: Otherwise → sold != true
									bson.M{"$ne": bson.A{"$$prop.sold", true}},
								},
							},
						},
					},
				},
			},
		}}},

		// Stage 4: Sort
		{{Key: "$sort", Value: bson.M{"_id": -1}}},

		// Stage 5: Skip
		{{Key: "$skip", Value: int64(skip)}},

		// Stage 6: Limit
		{{Key: "$limit", Value: int64(limit)}},
	}

	// ← ADD projection if fields are specified
	if len(fields) > 0 {
		projection := bson.M{}
		for _, field := range fields {
			projection[field] = 1
		}
		projection["_id"] = 1
		pipeline = append(pipeline, bson.D{{Key: "$project", Value: projection}})
	}

	// ← EXECUTE aggregation
	cursor, err := s.LeadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leads []models.Lead
	if err = cursor.All(ctx, &leads); err != nil {
		return nil, err
	}

	return leads, nil
}

func (s *LeadService) GetLeadPropertyDetails(ctx context.Context, leadID primitive.ObjectID) ([]bson.M, error) {
	// ← BUILD the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match the specific lead
		{{Key: "$match", Value: bson.M{"_id": leadID}}},

		// Stage 2: Lookup properties
		{{Key: "$lookup", Value: bson.M{
			"from":         "property",
			"localField":   "properties.property_id",
			"foreignField": "_id",
			"as":           "populated_properties",
		}}},

		// Stage 3: Add fields to combine property details with interest status
		{{Key: "$addFields", Value: bson.M{
			"properties_with_status": bson.M{
				"$map": bson.M{
					"input": "$properties",
					"as":    "prop",
					"in": bson.M{
						"$mergeObjects": bson.A{
							"$$prop",
							bson.M{
								"$arrayElemAt": bson.A{
									bson.M{
										"$filter": bson.M{
											"input": "$populated_properties",
											"cond": bson.M{
												"$eq": bson.A{"$$this._id", "$$prop.property_id"},
											},
										},
									},
									0,
								},
							},
						},
					},
				},
			},
		}}},

		// Stage 4: Filter out deleted properties
		{{Key: "$addFields", Value: bson.M{
			"properties_with_status": bson.M{
				"$filter": bson.M{
					"input": "$properties_with_status",
					"cond": bson.M{
						"$ne": []interface{}{"$$this.is_deleted", true},
					},
				},
			},
		}}},
		{{Key: "$addFields", Value: bson.M{
			"properties_with_status": bson.M{
				"$sortArray": bson.M{
					"input": "$properties_with_status",
					"sortBy": bson.M{
						"created_at": -1, // ← Sort by ObjectID descending (latest first)
					},
				},
			},
		}}},

		// Stage 5: Project only the properties array
		{{Key: "$project", Value: bson.M{
			"_id":        0,
			"properties": "$properties_with_status",
		}}},
	}

	// ← EXECUTE the aggregation pipeline
	cursor, err := s.LeadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// ← DECODE directly into bson.M array
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	// ← EXTRACT properties from first result
	if len(result) > 0 {
		if properties, ok := result[0]["properties"].(bson.A); ok {
			// Convert bson.A to []bson.M
			var propertiesArray []bson.M
			for _, prop := range properties {
				if propMap, ok := prop.(bson.M); ok {
					propertiesArray = append(propertiesArray, propMap)
				}
			}
			return propertiesArray, nil
		}
	}

	return []bson.M{}, nil
}

func (s *LeadService) GetConflictingProperties(ctx context.Context) ([]bson.M, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "property",
			"localField":   "properties.property_id",
			"foreignField": "_id",
			"as":           "populated_properties",
		}}},
		{{Key: "$unwind", Value: "$properties"}},
		{{Key: "$unwind", Value: "$populated_properties"}},
		{{Key: "$match", Value: bson.M{
			"$expr": bson.M{
				"$eq": []string{
					"$properties.property_id",
					"$populated_properties._id",
				},
			},
		}}},
		{{Key: "$match", Value: bson.M{
			"populated_properties.is_deleted": true,
			"populated_properties.sold":       bson.M{"$ne": true},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "dealers",
			"localField":   "properties.dealer_id",
			"foreignField": "_id",
			"as":           "dealer_info",
		}}},

		// Stage 6: Unwind dealer array (should be single item)
		{{Key: "$unwind", Value: "$dealer_info"}},
	}

	cursor, err := s.LeadCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

	return result, nil
}

func (s *LeadService) UpdatePropertyStatusByID(ctx context.Context, leadID primitive.ObjectID, propertyID primitive.ObjectID, status string) error {
	if status == "closed" {
		// Remove the property from the lead's properties array
		_, err := s.LeadCollection.UpdateOne(ctx,
			bson.M{"_id": leadID},
			bson.M{"$pull": bson.M{"properties": bson.M{"property_id": propertyID}}})
		return err
	} else {
		// Update the status
		_, err := s.LeadCollection.UpdateOne(ctx,
			bson.M{"_id": leadID, "properties.property_id": propertyID},
			bson.M{"$set": bson.M{"properties.$.status": status}})
		return err
	}
}
