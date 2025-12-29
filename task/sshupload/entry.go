package sshupload

import (
	"encoding/json"
	"fmt"
	"go-task/util"
	"os"
	"time"
)

type UploadInput struct {
	Config    []string `json:"config"`
	FileURL   string   `json:"file_url"`
	RemoteDir string   `json:"remote_dir"`
}

func Run(data string, sendFunc func(string)) (string, error) {
	var input UploadInput
	if err := json.Unmarshal([]byte(data), &input); err != nil {
		return "", fmt.Errorf("unmarshal input error: %s", err.Error())
	}
	if len(input.RemoteDir) == 0 {
		input.RemoteDir = "/tmp/task-sshupload-" + time.Now().Format("20060102150405")
	}

	tempFiles := util.DefaultSplitString(input.FileURL)
	if len(tempFiles) == 0 {
		sendFunc(fmt.Sprintf("file url error: %s", input.FileURL))
		return "", fmt.Errorf("file url error: %s", input.FileURL)
	}

	tempDir, err := os.MkdirTemp("", "task-sshupload-")
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

	retData := [][]string{}
	for _, item := range input.Config {
		sendFunc(fmt.Sprintf("upload to %s", item))
		remoteFiles, err := util.SSHUpload(item, localFiles, input.RemoteDir, sendFunc)
		if err != nil {
			sendFunc(fmt.Sprintf("upload error: %s", err.Error()))
			continue
		}

		retData = append(retData, remoteFiles)
	}

	byteData, err := json.Marshal(retData)
	if err != nil {
		sendFunc(fmt.Sprintf("marshal error: %s", err.Error()))
		return "", fmt.Errorf("marshal error: %s", err.Error())
	}

	return string(byteData), nil
}
