package md5

import (
	"encoding/json"
	"fmt"
	"go-task/util"
	"os"
)

// Md5Task 实现了Tool接口的演示任务
type Md5Task struct {
	Input   string `json:"input"`
	FileURL string `json:"file_url"`
}

// Run 执行任务，通过sendFunc发送结果
func Run(input string, sendFunc func(string)) (string, error) {
	var d Md5Task
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return "", err
	}
	if len(d.FileURL) > 0 {
		// 下载文件
		tempFiles := util.DefaultSplitString(d.FileURL)
		if len(tempFiles) == 0 {
			sendFunc(fmt.Sprintf("file url error: %s", d.FileURL))
			return "", fmt.Errorf("file url error: %s", d.FileURL)
		}

		tempDir, err := os.MkdirTemp("", "task-md5-")
		if err != nil {
			sendFunc(fmt.Sprintf("mkdir error: %s", err.Error()))
			return "", fmt.Errorf("mkdir error: %s", err.Error())
		}
		sendFunc(fmt.Sprintf("download %d files to %s", len(tempFiles), tempDir))
		localFiles, err := util.DownloadFiles(tempFiles, tempDir, sendFunc)
		if err != nil {
			sendFunc(fmt.Sprintf("download error: %s", err.Error()))
			return "", fmt.Errorf("download error: %s", err.Error())
		}

		// 计算文件MD5值
		for _, item := range localFiles {
			md5Value, err := util.CalculateFileMD5(item)
			if err != nil {
				sendFunc(fmt.Sprintf("calculate md5 error: %s", err.Error()))
				continue
			}
			sendFunc(fmt.Sprintf("file %s md5: %s", item, md5Value))
		}
	}

	if len(d.Input) > 0 {
		// 如果有输入，处理输入
		sendFunc(fmt.Sprintf("输入: %s", d.Input))
		// 计算MD5值
		md5Value := util.CalculateMD5(d.Input)
		sendFunc(fmt.Sprintf("MD5值: %s", md5Value))
	}

	return "", nil
}
