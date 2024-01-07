package main

import "fmt"

func getNestedField(data map[string]any, keys ...string) (any, error) {
	current := data

	level := 0
	for _, key := range keys {
		value, ok := current[key]
		if !ok {
			return nil, fmt.Errorf("key '%s' not found in the map", key)
		}

		// Check if the value is a map to continue the nesting
		if nested, isMap := value.(map[string]interface{}); isMap {
			current = nested
		} else {
			// If the value is not a map, it means we've reached the end of the nesting
			// and we return the final value
			if level < len(keys)-1 {
				return nil, fmt.Errorf("insufficient number of keys")
			}
			return value, nil
		}
		level++
	}

	return current, nil
}

func getNestedFieldWithDefault(defaultValue any, data map[string]interface{}, keys ...string) interface{} {
	res, err := getNestedField(data, keys...)
	if err != nil {
		return defaultValue
	}
	return res
}
