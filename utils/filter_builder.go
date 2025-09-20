package utils

import (
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BuildMongoFilter(params interface{}) bson.M {
	mongoFilter := bson.M{}

	
	arrayFieldsInMatch := getArrayFieldsInMatch(params)

	v := reflect.ValueOf(params)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			continue
		}

		fieldValue := v.Field(i)

		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		queryTag := field.Tag.Get("query")
		mongoTag := field.Tag.Get("mongo")
		convertTag := field.Tag.Get("convert")
		arrayTag := field.Tag.Get("array")
		operatorTag := field.Tag.Get("operator")

		if queryTag == "" {
			continue
		}

		var value interface{}
		if fieldValue.Kind() == reflect.Ptr {
			value = fieldValue.Elem().Interface()
		} else {
			value = fieldValue.Interface()
		}

		if convertTag != "" {
			value = applyMongoConversion(value, convertTag)
		}

		
		if boolValue, ok := value.(bool); ok && !boolValue {
			operatorTag = "$ne"
			value = true
		}

		fieldName := queryTag
		if mongoTag != "" {
			fieldName = mongoTag
		}

		if arrayTag != "" {
			if Contains(arrayFieldsInMatch, queryTag) {

				handleArrayFieldWithOperator(mongoFilter, fieldName, operatorTag, value)
			}

		} else {

			if operatorTag != "" {
				mongoFilter[fieldName] = bson.M{operatorTag: value}
				
			} else {
				mongoFilter[fieldName] = value
			}
		}
	}

	return mongoFilter
}

func getArrayFieldsInMatch(params interface{}) []string {
	v := reflect.ValueOf(params)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			if field.Type.Name() == "BaseQueryParams" {
				fieldValue := v.Field(i)
				if fieldValue.Kind() == reflect.Struct {
					arrayFieldsField := fieldValue.FieldByName("ArrayFilters")
					if arrayFieldsField.IsValid() && !arrayFieldsField.IsNil() {
						arrayFields := arrayFieldsField.Interface().(*string)
						if arrayFields != nil && *arrayFields != "" {
							fields := strings.Split(*arrayFields, ",")
							for i, field := range fields {
								fields[i] = strings.TrimSpace(field)
							}
							return fields
						}
					}
				}
			}
		}
	}

	return []string{} // ✅ Default: no array fields in match
}



// ✅ Generic array field handler
func handleArrayFieldWithOperator(mongoFilter bson.M, fieldName string, operator string, value interface{}) {
	parts := strings.Split(fieldName, ".")
	if len(parts) != 2 {
		if operator != "" {
			mongoFilter[fieldName] = bson.M{operator: value}
		} else {
			mongoFilter[fieldName] = value
		}
		return
	}

	arrayField := parts[0]
	nestedField := parts[1]

	
	

	switch operator {
	case "$ne", "$nin":
		mongoFilter[arrayField] = bson.M{"$not": bson.M{"$elemMatch": bson.M{nestedField: value}}}
	case "$all":
		mongoFilter[arrayField] = bson.M{"$all": value}
	case "$not":
		mongoFilter[arrayField] = bson.M{"$not": value}
	case "$exists", "$size", "$type":
		mongoFilter[arrayField] = bson.M{"$elemMatch": bson.M{nestedField: bson.M{operator: value}}}
	case "":
		mongoFilter[arrayField] = bson.M{"$elemMatch": bson.M{nestedField: value}}
	default:
		mongoFilter[arrayField] = bson.M{"$elemMatch": bson.M{nestedField: bson.M{operator: value}}}
	}
}

func NeedsAggregation(filter bson.M) bool {
	for _, value := range filter {
		if elemMatch, ok := value.(bson.M); ok {

			if _, hasElemMatch := elemMatch["$elemMatch"]; hasElemMatch {
				return true
			}

			arrayOperators := []string{"$not", "$all", "$size", "$exists", "$type", "$regex"}
			for _, op := range arrayOperators {
				if _, hasOp := elemMatch[op]; hasOp {
					return true
				}
			}
		}
	}
	return false
}

func applyMongoConversion(value interface{}, convertType string) interface{} {
	switch convertType {
	case "objectid":
		if str, ok := value.(string); ok {
			if objectID, err := primitive.ObjectIDFromHex(str); err == nil {
				return objectID
			}
		}
	case "date":
		if str, ok := value.(string); ok {
			if date, err := time.Parse("2006-01-02", str); err == nil {
				return date
			}
		}
	}
	return value
}

func BuildAggregationPipeline(filter bson.M, sortField string, sortOrder int, skip int64, limit int64, projection bson.M) []bson.M {
	pipeline := []bson.M{
		{"$match": filter},
	}

	arrayFieldGroups := make(map[string]bson.M)

	for fieldName, fieldValue := range filter {
		if elemMatch, ok := fieldValue.(bson.M); ok {
			if _, hasElemMatch := elemMatch["$elemMatch"]; hasElemMatch {
				condition := elemMatch["$elemMatch"].(bson.M)

				if arrayFieldGroups[fieldName] == nil {
					arrayFieldGroups[fieldName] = bson.M{}
				}

				for nestedField, nestedValue := range condition {
					arrayFieldGroups[fieldName][nestedField] = nestedValue
				}
			}
		}
	}

	for fieldName, combinedCondition := range arrayFieldGroups {
		var andConditions []bson.M

		for nestedField, nestedValue := range combinedCondition {
			// ✅ Check if nestedValue is already a MongoDB operator structure
			if operatorMap, ok := nestedValue.(bson.M); ok {
				// ✅ Extract the operator and value
				for operator, value := range operatorMap {
					andConditions = append(andConditions, bson.M{
						operator: bson.A{"$$item." + nestedField, value},
					})
				}
			} else {
				// ✅ Plain value: use $eq
				andConditions = append(andConditions, bson.M{
					"$eq": bson.A{"$$item." + nestedField, nestedValue},
				})
			}
		}

		var condExpr interface{}
		if len(andConditions) == 1 {
			condExpr = andConditions[0]
		} else {
			condExpr = bson.M{"$and": andConditions}
		}

		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				fieldName: bson.M{
					"$filter": bson.M{
						"input": "$" + fieldName,
						"as":    "item",
						"cond":  condExpr,
					},
				},
			},
		})
	}

	if sortField != "" {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{sortField: sortOrder}})
	}

	if skip > 0 {
		pipeline = append(pipeline, bson.M{"$skip": skip})
	}
	if limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}

	if projection != nil {
		pipeline = append(pipeline, bson.M{"$project": projection})
	}

	return pipeline
}
