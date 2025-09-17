// utils/field_selector.go
package utils

import (
    "net/http"
    "strings"
    "go.mongodb.org/mongo-driver/bson"
)


func ParseFieldSelection(r *http.Request) []string {
    queryParams := r.URL.Query()
    
    if fieldsParam := queryParams.Get("fields"); fieldsParam != "" {
        return strings.Split(fieldsParam, ",")
    }
    
    return []string{} 
}


func BuildMongoProjection(fields []string) bson.M {
    if len(fields) == 0 {
        return bson.M{} 
    }
    
    projection := bson.M{}
    for _, field := range fields {
        projection[field] = 1
    }
    
    
    projection["_id"] = 1
    
    return projection
}


