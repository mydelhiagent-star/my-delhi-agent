// utils/response_filter.go
package utils

import (
    "reflect"
    "strings"
)

// ✅ Generic response field filtering
func FilterResponseFields(data interface{}, fields []string) interface{} {
    if len(fields) == 0 {
        return data
    }
    
    v := reflect.ValueOf(data)
    t := v.Type()
    
    // ✅ Handle slice
    if t.Kind() == reflect.Slice {
        filteredSlice := reflect.MakeSlice(t, 0, v.Len())
        
        for i := 0; i < v.Len(); i++ {
            filteredItem := filterStructFields(v.Index(i), fields)
            filteredSlice = reflect.Append(filteredSlice, filteredItem)
        }
        
        return filteredSlice.Interface()
    }
    
    // ✅ Handle single struct
    if t.Kind() == reflect.Struct {
        return filterStructFields(v, fields)
    }
    
    return data
}

// ✅ Filter struct fields
func filterStructFields(v reflect.Value, fields []string) reflect.Value {
    t := v.Type()
    filteredStruct := reflect.New(t).Elem()
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fieldValue := v.Field(i)
        
        // ✅ Check if field is requested
        if contains(fields, getFieldName(field)) {
            filteredStruct.Field(i).Set(fieldValue)
        }
    }
    
    return filteredStruct
}

// ✅ Get field name from struct tag
func getFieldName(field reflect.StructField) string {
    jsonTag := field.Tag.Get("json")
    if jsonTag != "" {
        return strings.Split(jsonTag, ",")[0]
    }
    return strings.ToLower(field.Name)
}

// ✅ Check if field is in list
func contains(fields []string, field string) bool {
    for _, f := range fields {
        if f == field {
            return true
        }
    }
    return false
}