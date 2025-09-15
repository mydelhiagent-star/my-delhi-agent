// utils/filter_builder.go
package utils

import (
    "reflect"
    "strings"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func BuildFilters(params interface{}) map[string]interface{} {
    filters := make(map[string]interface{})
    
    v := reflect.ValueOf(params)
    t := v.Type()
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fieldValue := v.Field(i)
        
        if fieldValue.IsNil() {
            continue
        }
        
        queryTag := field.Tag.Get("query")
        if queryTag == "" {
            continue
        }
        
        if shouldIncludeField(fieldValue, queryTag) {
            value := fieldValue.Elem().Interface()
            
            // ✅ Smart ID detection - no hardcoding
            if isIDField(queryTag) {
                value = convertToObjectID(value)
            }

			mongoFieldName := getMongoFieldName(queryTag)
            filters[mongoFieldName] = value

			
        }
    }
    
    return filters
}

// ✅ Generic ID detection
func isIDField(fieldName string) bool {
    // Check if field name ends with "_id" or is "id"
    return fieldName == "id" || strings.HasSuffix(fieldName, "_id")
}

func getMongoFieldName(queryTag string) string {
    if queryTag == "id" {
        return "_id"  // Frontend "id" -> MongoDB "_id"
    }
    return queryTag  // Other fields stay the same
}

// ✅ Generic ObjectID conversion
func convertToObjectID(value interface{}) interface{} {
    if str, ok := value.(string); ok {
        if objectID, err := primitive.ObjectIDFromHex(str); err == nil {
            return objectID
        }
    }
    return value // Return original if conversion fails
}

func shouldIncludeField(fieldValue reflect.Value, queryTag string) bool {
    if queryTag == "page" || queryTag == "limit" {
        return false
    }
    return true
}