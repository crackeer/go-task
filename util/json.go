package util

import "encoding/json"

func ParseJSON(value string) interface{} {
	var result interface{}
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return nil
	}
	return result
}
