// utils/filter_builder.go
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
        
        
        if isNestedArrayField(fieldName) {
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