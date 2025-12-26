package sshupload

import (
	"encoding/json"
	"fmt"
	"go-task/util"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

	tempFiles := strings.Split(input.FileURL, ",")
	if len(tempFiles) == 0 {
		sendFunc(fmt.Sprintf("file url error: %s", input.FileURL))
		return "", fmt.Errorf("file url error: %s", input.FileURL)
	}

	var files []string
	for _, file := range tempFiles {
		files = append(files, strings.Split(strings.TrimSpace(file), "\n")...)
	}

	retData := [][]string{}
	for _, item := range input.Config {
		sendFunc(fmt.Sprintf("upload to %s", item))
		remoteFiles, err := doUpload(item, files, input.RemoteDir, sendFunc)
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

func doUpload(config string, fileURLs []string, remoteDir string, sendFunc func(string)) ([]string, error) {
	remoteFiles := make([]string, 0)
	parts := strings.Split(config, " ")
	if len(parts) != 4 {
		sendFunc(fmt.Sprintf("config error: %s", config))
		return nil, fmt.Errorf("config error: %s", config)
	}
	host := parts[0]
	port := parts[1]
	user := parts[2]
	password := parts[3]
	sshClient, err := util.NewSSHClient(host, port, user, password)
	if err != nil {
		sendFunc(fmt.Sprintf("new ssh client error: %s", err.Error()))
		return nil, fmt.Errorf("new ssh client error: %s", err.Error())
	}
	defer sshClient.Close()

	sendFunc(fmt.Sprintf("remote mkdir %s", remoteDir))
	if err := sshClient.Mkdir(remoteDir); err != nil {
		sendFunc(fmt.Sprintf("mkdir error: %s", err.Error()))
		return nil, fmt.Errorf("mkdir error: %s", err.Error())
	}

	for _, fileURL := range fileURLs {
		tempDir := os.TempDir()
		urlObject, err := url.Parse(fileURL)
		if err != nil {
			sendFunc(fmt.Sprintf("parse url error: %s", err.Error()))
			continue
		}
		targetFile := filepath.Join(tempDir, filepath.Base(urlObject.Path))
		sendFunc(fmt.Sprintf("download file %s to %s...", fileURL, targetFile))
		if _, err := util.DownloadToDest(fileURL, targetFile); err != nil {
			sendFunc(fmt.Sprintf("download error: %s", err.Error()))
			continue
		}
		sendFunc(fmt.Sprintf("download file %s to %s success", fileURL, targetFile))
		remoteFile := filepath.Join(remoteDir, filepath.Base(urlObject.Path))

		// 2. upload file
		sendFunc(fmt.Sprintf("upload file %s to %s....", targetFile, remoteFile))
		if err := sshClient.UploadTo(targetFile, remoteFile); err != nil {
			sendFunc(fmt.Sprintf("upload error: %s", err.Error()))
			continue
		}
		sendFunc(fmt.Sprintf("upload file %s to %s success", targetFile, remoteFile))
		remoteFiles = append(remoteFiles, remoteFile)
	}
	return remoteFiles, nil
}
