package json2csv

import (
	"encoding/json"
	"fmt"
	"go-task/util"
	"os"
)

// Json2CsvTask 实现了Tool接口的演示任务
type Json2CsvTask struct {
	Input   string `json:"input"`
	TempDir string `json:"tempDir"`
}

// Run 执行任务，通过sendFunc发送结果
func Run(input string, sendFunc func(string)) (string, error) {
	var d Json2CsvTask
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return "", err
	}
	// 转换JSON为CSV
	sendFunc("正在转换JSON为CSV...")
	csvFilePath := fmt.Sprintf("%s/output.csv", d.TempDir)
	if err := os.MkdirAll(d.TempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}
	if err := util.Json2CsvFile([]byte(d.Input), csvFilePath); err != nil {
		return "", fmt.Errorf("failed to convert JSON to CSV: %v", err)
	}

	sendFunc(fmt.Sprintf("CSV文件已生成: %s", csvFilePath))
	downloadURL := util.GenerationDownloadURL(csvFilePath)
	sendFunc("下载CSV：" + downloadURL)

	return downloadURL, nil
}
