package utils

import (
    "net/http"
    "net/url"
    "reflect"
    "strconv"
)

func ParseQueryParams(r *http.Request, params interface{}) error {
    queryParams := r.URL.Query()
    v := reflect.ValueOf(params).Elem()
    t := v.Type()
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fieldValue := v.Field(i)
        
        // ✅ Handle embedded structs (Google's approach)
        if field.Anonymous {
            if err := parseEmbeddedStruct(fieldValue, queryParams); err != nil {
                return err
            }
            continue
        }
        
        queryTag := field.Tag.Get("query")
        if queryTag == "" {
            continue
        }
        
        value := queryParams.Get(queryTag)
        if value == "" {
            continue
        }
        
        setPointerValue(fieldValue, value)
    }
    
    return nil
}

// ✅ Google's pattern: Recursive parsing for embedded structs
func parseEmbeddedStruct(structValue reflect.Value, queryParams url.Values) error {
    if structValue.Kind() != reflect.Struct {
        return nil
    }
    
    structType := structValue.Type()
    for i := 0; i < structType.NumField(); i++ {
        field := structType.Field(i)
        fieldValue := structValue.Field(i)
        
        queryTag := field.Tag.Get("query")
        if queryTag == "" {
            continue
        }
        
        value := queryParams.Get(queryTag)
        if value == "" {
            continue
        }
        
        setPointerValue(fieldValue, value)
    }
    
    return nil
}

func setPointerValue(fieldValue reflect.Value, value string) {
    elemType := fieldValue.Type().Elem()
    
    switch elemType.Kind() {
    case reflect.String:
        // ✅ Always set string value (even if empty)
        fieldValue.Set(reflect.ValueOf(&value))
        
    case reflect.Bool:
        if boolVal, err := strconv.ParseBool(value); err == nil {
            fieldValue.Set(reflect.ValueOf(&boolVal))
        }
        
    case reflect.Int:
        if intVal, err := strconv.Atoi(value); err == nil {
            fieldValue.Set(reflect.ValueOf(&intVal))
        }
        
    case reflect.Float64:
        if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
            fieldValue.Set(reflect.ValueOf(&floatVal))
        }
    }
}