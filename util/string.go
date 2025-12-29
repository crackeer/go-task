package util

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"strings"
)

// CalculateMD5 计算字符串的MD5值
func CalculateMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return fmt.Sprintf("%x", hash)
}

func GenerationDownloadURL(filePath string) string {
	return fmt.Sprintf("http://{WINDOW_HOSTNAME}/api/download?file_path=%s", filePath)
}

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Base64Decode(input string) string {
	decoded, _ := base64.StdEncoding.DecodeString(input)
	return string(decoded)
}

func DefaultSplitString(input string) []string {
	return splitString([]string{input}, []string{" ", "\n", "\t", ",", ";"})
}

func splitString(input []string, delimiter []string) []string {
	if len(delimiter) < 1 {
		return input
	}
	retData := make([]string, 0)
	first := delimiter[0]

	for _, item := range input {
		parts := strings.Split(item, first)
		retData = append(retData, parts...)
	}
	return splitString(retData, delimiter[1:])
}
