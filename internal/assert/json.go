package assert

import (
	"encoding/json"
	"fmt"
	"strings"
)

func extractvalue(data interface{}, path string) (interface{}, error) {
	if path == "$.length" {
		if arr, ok := data.([]interface{}); ok {
			return len(arr), nil
		}
		return nil, fmt.Errorf("$.length applied to non-array")
	}

	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		obj, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid path: %s", path)
		}

		val, exists := obj[part]
		if !exists {
			return nil, fmt.Errorf("field %s not found", path)
		}
		current = val
	}

	return current, nil
}

func isComparison(v string) bool {
	return strings.HasPrefix(v, ">") || strings.HasPrefix(v, "<")
}

func compare(actual interface{}, expStr string) error {
	var actualNum float64

	switch v := actual.(type) {
	case float64:
		actualNum = v
	case int:
		actualNum = float64(v)
	case int64:
		actualNum = float64(v)
	default:
		return fmt.Errorf("comparison on non-number")
	}

	rule := strings.TrimSpace(expStr)
	operator := rule[:1]
	value := rule[1:]
	var expected float64
	fmt.Sscanf(value, "%f", &expected)

	switch operator {
	case ">":
		if actualNum <= expected {
			return fmt.Errorf("%v is not > %v", actualNum, expected)
		}
	case "<":
		if actualNum >= expected {
			return fmt.Errorf("%v is not < %v", actualNum, expected)
		}
	}

	return nil

}

func AssertJSON(body []byte, rules map[string]interface{}) error {
	var data interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("invalid json response")
	}
	for path, expected := range rules {
		actual, err := extractvalue(data, path)
		if err != nil {
			return err
		}

		if expStr, ok := expected.(string); ok && isComparison(expStr) {
			if err := compare(actual, expStr); err != nil {
				return fmt.Errorf("assertion failed at %s: %v", path, err)
			}

			continue
		}

		if fmt.Sprint(actual) != fmt.Sprint(expected) {
			return fmt.Errorf("assertion failed at %s: expected %v, got %v", path, expected, actual)
		}
	}
	return nil
}
