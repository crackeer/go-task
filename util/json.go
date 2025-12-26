package util

import (
	"encoding/json"
	"net/url"

	"github.com/tidwall/gjson"
)

func ParseJSON(value string) interface{} {
	var result interface{}
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return nil
	}
	return result
}

// IsURL
//
//	@param value
//	@return bool
func IsURL(value string) bool {
	parsedURL, err := url.Parse(value)
	return err == nil && parsedURL.Host != ""
}

// ExtractURL
//
//	@param value
//	@param matchFunc
//	@return []string
func ExtractURL(value interface{}, matchFunc func(string) bool) []string {
	if strValue, ok := value.(string); ok {
		if gjson.Valid(strValue) {
			var data interface{}
			if err := json.Unmarshal([]byte(strValue), &data); err == nil {
				return ExtractURL(data, matchFunc)
			}
		}
		if matchFunc(strValue) {
			return []string{strValue}
		}
		return []string{}
	}

	var retData []string
	if mapValue, ok := value.(map[string]interface{}); ok {
		for _, value := range mapValue {
			retData = append(retData, ExtractURL(value, matchFunc)...)
		}
		return retData
	}

	if listValue, ok := value.([]interface{}); ok {
		for _, value := range listValue {
			retData = append(retData, ExtractURL(value, matchFunc)...)
		}
		return retData
	}
	return nil
}
