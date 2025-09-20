// utils/update_builder.go
package utils

import (
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func BuildUpdateDocument(update interface{}) bson.M {
    updateDoc := bson.M{}
    
    v := reflect.ValueOf(update)
    t := reflect.TypeOf(update)
    
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        fieldValue := v.Field(i)
        
        // Skip if field is nil
        if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
            continue
        }
        
       
        bsonTag := field.Tag.Get("bson")
        fieldName := bsonTag  
       
        
        // Add to update document
        if fieldValue.Kind() == reflect.Ptr {
            updateDoc[fieldName] = fieldValue.Elem().Interface()
        } else {
            updateDoc[fieldName] = fieldValue.Interface()
        }
    }

	updateDoc["updated_at"] = time.Now()
    
    return updateDoc
}