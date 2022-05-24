package myjson

import (
	"encoding/json"
	"fmt"
)

// Encode encode value into JSON string
func Encode(v interface{}) string {
	str, err := json.Marshal(v)

	if err != nil {
		fmt.Printf("[ERROR] encode value to JSON: %v", err)
	}

	return string(str)
}

// EncodePretty encode value into JSON string
func EncodePretty(v interface{}) string {
	str, err := json.MarshalIndent(v, "", "\t")

	if err != nil {
		fmt.Printf("[ERROR] encode value to JSON: %v", err)
	}

	return string(str)
}

// Decode decode and parse JSON string into struct
func Decode(jsonText string, holder interface{}) {
	err := json.Unmarshal([]byte(jsonText), &holder)
	if err != nil {
		fmt.Printf("[ERROR] decode JSON value: %v", err)
	}
}
