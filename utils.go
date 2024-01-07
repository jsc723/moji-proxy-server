package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

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


func loadConfig(path string) Config {
	var c Config
	// Open and read the configuration file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening configuration file:", err)
		return c
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		return c
	}

	err = json.Unmarshal(content, &c)
	if err != nil {
		fmt.Println("Error parsing configuration file:", err)
		c = Config{}
		return c
	}
	return c
}