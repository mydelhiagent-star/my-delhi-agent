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

	// ✅ Get array fields that should be in match
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
		operatorTag := field.Tag.Get("operator")
		arrayTag := field.Tag.Get("array")

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

		// ✅ Handle boolean logic: false -> $ne: true
		if boolValue, ok := value.(bool); ok && !boolValue {
			operatorTag = "ne"
			value = true
		}

		fieldName := queryTag
		if mongoTag != "" {
			fieldName = mongoTag
		}

		// ✅ Check if array field should be in match
		if arrayTag != "" {
			if Contains(arrayFieldsInMatch, queryTag) {
				// ✅ Include array field in match
				handleArrayFieldWithOperator(mongoFilter, fieldName, operatorTag, value)
			}
			// ✅ If not in arrayFieldsInMatch, skip from match (will go to addFields)
		} else {
			// ✅ Regular fields always go to match
			if operatorTag != "" {
				mongoOperator := mapOperatorToMongo(operatorTag)
				if mongoOperator != "" {
					mongoFilter[fieldName] = bson.M{mongoOperator: value}
				}
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

func mapOperatorToMongo(operator string) string {
	operatorMap := map[string]string{
		"eq":     "$eq",
		"ne":     "$ne",
		"nin":    "$nin",
		"in":     "$in",
		"gte":    "$gte",
		"lte":    "$lte",
		"gt":     "$gt",
		"lt":     "$lt",
		"regex":  "$regex",
		"exists": "$exists",
		"size":   "$size",
	}
	return operatorMap[operator]
}

// ✅ Generic array field handler
func handleArrayFieldWithOperator(mongoFilter bson.M, fieldName string, operator string, value interface{}) {
	parts := strings.Split(fieldName, ".")
	if len(parts) != 2 {
		mongoOperator := mapOperatorToMongo(operator)
		if mongoOperator != "" {
			mongoFilter[fieldName] = bson.M{mongoOperator: value}
		} else {
			mongoFilter[fieldName] = value
		}
		return
	}

	arrayField := parts[0]
	nestedField := parts[1]

	// ✅ Convert operator to MongoDB operator
	mongoOperator := mapOperatorToMongo(operator)

	// ✅ Generic logic for array field operators
	switch mongoOperator {
	case "$ne", "$nin":
		// ✅ Exclusion operators: Use $not with $elemMatch
		mongoFilter[arrayField] = bson.M{
			"$not": bson.M{
				"$elemMatch": bson.M{
					nestedField: value,
				},
			},
		}
	case "":
		// ✅ No operator: Direct equality
		mongoFilter[arrayField] = bson.M{
			"$elemMatch": bson.M{
				nestedField: value,
			},
		}
	default:
		// ✅ Other operators: Use $elemMatch with operator
		mongoFilter[arrayField] = bson.M{
			"$elemMatch": bson.M{
				nestedField: bson.M{
					mongoOperator: value,
				},
			},
		}
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
