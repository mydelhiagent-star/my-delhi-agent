package utils

import (
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// ✅ Update utils/filter_builder.go (replace lines 55-67)
func BuildMongoFilter(params interface{}) bson.M {
    mongoFilter := bson.M{}
    
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
        operatorTag := field.Tag.Get("operator") // ✅ Use operator tag
        
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
        
        fieldName := queryTag
        if mongoTag != "" {
            fieldName = mongoTag
        }
        

        if operatorTag != "" {
            mongoOperator := mapOperatorToMongo(operatorTag)
            
            if mongoOperator != "" {
                
                if isNestedArrayField(fieldName) {
                    
                    handleArrayFieldWithOperator(mongoFilter, fieldName, mongoOperator, value)
                } else {
                   
                    mongoFilter[fieldName] = bson.M{mongoOperator: value}
                }
            }
        } else if isNestedArrayField(fieldName) {
           
            handleNestedArrayField(mongoFilter, fieldName, value)
        } else {
            
            if boolValue, ok := value.(bool); ok {
                if !boolValue {
                    mongoFilter[fieldName] = bson.M{"$ne": true}
                } else {
                    mongoFilter[fieldName] = true
                }
            } else {
                mongoFilter[fieldName] = value
            }
        }
    }
    
    return mongoFilter
}


func mapOperatorToMongo(operator string) string {
    operatorMap := map[string]string{
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
        mongoFilter[fieldName] = bson.M{operator: value}
        return
    }
    
    arrayField := parts[0]
    nestedField := parts[1]
    
    // ✅ Generic logic for array field operators
    switch operator {
    case "$ne", "$nin":
        // ✅ Exclusion operators: Use $not with $elemMatch
        mongoFilter[arrayField] = bson.M{
            "$not": bson.M{
                "$elemMatch": bson.M{
                    nestedField: value,
                },
            },
        }
    default:
        // ✅ Inclusion operators: Use $elemMatch
        mongoFilter[arrayField] = bson.M{
            "$elemMatch": bson.M{
                nestedField: bson.M{
                    operator: value,
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

func isNestedArrayField(fieldName string) bool {
    return strings.Contains(fieldName, ".")
}

func handleNestedArrayField(mongoFilter bson.M, fieldName string, value interface{}) {
    parts := strings.Split(fieldName, ".")
    if len(parts) != 2 {
        mongoFilter[fieldName] = value
        return
    }
    
    arrayField := parts[0]       
    nestedField := parts[1]       
    
    
    if existingElemMatch, exists := mongoFilter[arrayField]; exists {
        if elemMatch, ok := existingElemMatch.(bson.M); ok {
            if _, hasElemMatch := elemMatch["$elemMatch"]; hasElemMatch {
              
                elemMatch["$elemMatch"].(bson.M)[nestedField] = value
                return
            }
        }
    }
    
   
    mongoFilter[arrayField] = bson.M{
        "$elemMatch": bson.M{
            nestedField: value,
        },
    }
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
            andConditions = append(andConditions, bson.M{
                "$eq": bson.A{"$$item." + nestedField, nestedValue},
            })
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