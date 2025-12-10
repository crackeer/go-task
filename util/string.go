package util

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
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
