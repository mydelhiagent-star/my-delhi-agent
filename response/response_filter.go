package response

import (
    "encoding/json"
)





func filterDataFields(data interface{}, fields []string) interface{} {
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return data
    }
    
    var jsonMap interface{}
    if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
        return data
    }
    
    return filterJSONMap(jsonMap, fields)
}


func filterJSONMap(data interface{}, fields []string) interface{} {
    switch v := data.(type) {
    case map[string]interface{}:
        filtered := make(map[string]interface{})
        for key, value := range v {
            if contains(fields, key) {
                filtered[key] = value
            }
        }
        return filtered
        
    case []interface{}:
        var filtered []interface{}
        for _, item := range v {
            filtered = append(filtered, filterJSONMap(item, fields))
        }
        return filtered
        
    default:
        return data
    }
}


func contains(fields []string, field string) bool {
    for _, f := range fields {
        if f == field {
            return true
        }
    }
    return false
}