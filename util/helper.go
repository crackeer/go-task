package util

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DownloadFiles(files []string, tempDir string, sendFunc func(string)) ([]string, error) {
	sendFunc(fmt.Sprintf("download %d files", len(files)))
	localFiles := make([]string, 0)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		sendFunc(fmt.Sprintf("mkdir error: %s", err.Error()))
		return nil, fmt.Errorf("mkdir error: %s", err.Error())
	}
	for _, fileURL := range files {
		sendFunc(fmt.Sprintf("download file %s...", fileURL))
		urlObject, err := url.Parse(fileURL)
		if err != nil {
			sendFunc(fmt.Sprintf("parse url error: %s", err.Error()))
			continue
		}
		targetFile := filepath.Join(tempDir, filepath.Base(urlObject.Path))
		sendFunc(fmt.Sprintf("download file %s to %s...", fileURL, targetFile))
		if _, err := DownloadToDest(fileURL, targetFile); err != nil {
			sendFunc(fmt.Sprintf("download error: %s", err.Error()))
			continue
		}
		localFiles = append(localFiles, targetFile)
		sendFunc(fmt.Sprintf("download file %s to %s success", fileURL, targetFile))
	}
	return localFiles, nil
}

func SSHUpload(sshClient *SSHClient, localFiles []string, remoteDir string, sendFunc func(string)) ([]string, error) {
	remoteFiles := make([]string, 0)
	if len(localFiles) == 0 {
		sendFunc(fmt.Sprintf("local files is empty"))
		return nil, fmt.Errorf("local files is empty")
	}

	sendFunc(fmt.Sprintf("remote mkdir %s", remoteDir))
	if err := sshClient.Mkdir(remoteDir); err != nil {
		sendFunc(fmt.Sprintf("mkdir error: %s", err.Error()))
		return nil, fmt.Errorf("mkdir error: %s", err.Error())
	}
	sendFunc(fmt.Sprintf("local files: %s", strings.Join(localFiles, ",")))
	for _, localFile := range localFiles {
		remoteFile := filepath.Join(remoteDir, filepath.Base(localFile))
		// 2. upload file
		sendFunc(fmt.Sprintf("upload file %s to %s....", localFile, remoteFile))
		if err := sshClient.UploadTo(localFile, remoteFile); err != nil {
			sendFunc(fmt.Sprintf("upload error: %s", err.Error()))
			continue
		}
		sendFunc(fmt.Sprintf("upload file %s to %s success", localFile, remoteFile))
		remoteFiles = append(remoteFiles, remoteFile)
	}
	return remoteFiles, nil
}

func NewSSHClientByConfig(config string) (*SSHClient, error) {
	parts := strings.Split(config, " ")
	if len(parts) != 4 {
		return nil, fmt.Errorf("config error: %s", config)
	}
	host := parts[0]
	port := parts[1]
	user := parts[2]
	password := parts[3]
	return NewSSHClient(host, port, user, password)
}
