package json2csv

import (
	"encoding/json"
	"fmt"
	"go-task/util"
	"os"
	"path/filepath"
)

// Json2CsvTask 实现了Tool接口的演示任务
type Json2CsvTask struct {
	Input string `json:"input"`
}

// Run 执行任务，通过sendFunc发送结果
func Run(input string, sendFunc func(string)) (string, error) {
	var d Json2CsvTask
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return "", err
	}
	// 转换JSON为CSV
	sendFunc("正在转换JSON为CSV...")
	tempDir, err := os.MkdirTemp("", "json2csv-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}
	csvFilePath := filepath.Join(tempDir, "output.csv")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}
	if err := util.Json2CsvFile([]byte(d.Input), csvFilePath); err != nil {
		return "", fmt.Errorf("failed to convert JSON to CSV: %v", err)
	}

	sendFunc(fmt.Sprintf("CSV文件已生成: %s", csvFilePath))
	// 上传CSV文件
	sendFunc("正在上传CSV文件...")
	uploadResp, err := util.UploadFile(csvFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to upload CSV file: %v", err)
	}
	sendFunc(fmt.Sprintf("CSV文件已上传: %s", uploadResp.URL))

	return uploadResp.URL, nil
}
