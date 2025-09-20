// utils/update_builder.go
package utils

import (
	"reflect"
	"strings"

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
        
        // ✅ Get BSON field name directly
        bsonTag := field.Tag.Get("bson")
        fieldName := bsonTag  // ✅ No split needed!
        if fieldName == "" {
            fieldName = strings.ToLower(field.Name)
        }
        
        // Add to update document
        if fieldValue.Kind() == reflect.Ptr {
            updateDoc[fieldName] = fieldValue.Elem().Interface()
        } else {
            updateDoc[fieldName] = fieldValue.Interface()
        }
    }
    
    return updateDoc
}