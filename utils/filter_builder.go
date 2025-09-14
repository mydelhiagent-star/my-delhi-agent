package utils

import "reflect"

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
            filters[queryTag] = fieldValue.Elem().Interface()
        }
    }
    
    return filters
}

func shouldIncludeField(fieldValue reflect.Value, queryTag string) bool {
    if queryTag == "page" || queryTag == "limit" {
        return false
    }
    return true
}