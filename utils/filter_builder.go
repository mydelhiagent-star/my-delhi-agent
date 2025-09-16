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
        if field.Anonymous{
            continue
        }
        fieldValue := v.Field(i)
        if fieldValue.IsNil() {
            continue
        }
        
        queryTag := field.Tag.Get("query")
        mongoTag := field.Tag.Get("mongo")
        convertTag := field.Tag.Get("convert")
		paginationTag := field.Tag.Get("pagination")
        
        if queryTag == "" {
            continue
        }

		if paginationTag == "true"{
			continue
		}
        
        value := fieldValue.Elem().Interface()
        
        
        if convertTag != "" {
            value = applyMongoConversion(value, convertTag)
        }
        
        if mongoTag != "" {
            mongoFilter[mongoTag] = value
        } else {
            mongoFilter[queryTag] = value
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