// utils/field_selector.go
package utils

import (
    "net/http"
    "strings"
    "go.mongodb.org/mongo-driver/bson"
)

// ✅ Generic field selection for any model
type FieldSelector struct {
    Fields []string
}

// ✅ Parse field selection from URL
func ParseFieldSelection(r *http.Request) *FieldSelector {
    queryParams := r.URL.Query()
    
    selector := &FieldSelector{}
    
    if fieldsParam := queryParams.Get("fields"); fieldsParam != "" {
        selector.Fields = strings.Split(fieldsParam, ",")
    }
    
    return selector
}

// ✅ Build MongoDB projection for any model
func (fs *FieldSelector) BuildProjection(fieldMappings map[string]string) bson.M {
    if len(fs.Fields) == 0 {
        return bson.M{} // ✅ No projection - return all fields
    }
    
    projection := bson.M{}
    for _, field := range fs.Fields {
        mongoField := fs.mapFieldToMongo(field, fieldMappings)
        projection[mongoField] = 1
    }
    
    // ✅ Always include _id
    projection["_id"] = 1
    
    return projection
}

// ✅ Map domain field to MongoDB field
func (fs *FieldSelector) mapFieldToMongo(field string, fieldMappings map[string]string) string {
    if mongoField, exists := fieldMappings[field]; exists {
        return mongoField
    }
    return field
}