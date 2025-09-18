// utils/filter_builder.go
package utils

import (
	"reflect"
	
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
        
        // Get value safely
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
    
    return mongoFilter
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