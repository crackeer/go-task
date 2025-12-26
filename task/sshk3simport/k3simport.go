package sshk3simport

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

type K3sImportInput struct {
	Config     string `json:"config"`
	FileURL    string `json:"file_url"`
	Namespace  string `json:"namespace"`
	Deployment string `json:"deployment"`
}

// RunK3sImport
//
//	@param data
//	@param sendFunc
//	@return string
//	@return error
func RunK3sImport(data string, sendFunc func(string)) (string, error) {
	var input K3sImportInput
	if err := json.Unmarshal([]byte(data), &input); err != nil {
		return "", fmt.Errorf("unmarshal input error: %s", err.Error())
	}
	parts := strings.Split(input.Config, " ")
	if len(parts) != 4 {
		sendFunc(fmt.Sprintf("config error: %s", input.Config))
		return "", fmt.Errorf("config error: %s", input.Config)
	}
	host := parts[0]
	port := parts[1]
	user := parts[2]
	password := parts[3]
	sshClient, err := util.NewSSHClient(host, port, user, password)
	if err != nil {
		sendFunc(fmt.Sprintf("new ssh client error: %s", err.Error()))
		return "", fmt.Errorf("new ssh client error: %s", err.Error())
	}
	defer sshClient.Close()
	tempRemoteDir := "/tmp/task-sshupload-" + time.Now().Format("20060102150405")
	sendFunc(fmt.Sprintf("remote mkdir %s", tempRemoteDir))
	if err := sshClient.Mkdir(tempRemoteDir); err != nil {
		sendFunc(fmt.Sprintf("mkdir error: %s", err.Error()))
		return "", fmt.Errorf("mkdir error: %s", err.Error())
	}

	// 1. download file
	tempDir := os.TempDir()
	sendFunc(fmt.Sprintf("download file %s to %s", input.FileURL, tempDir))
	urlObject, err := url.Parse(input.FileURL)
	if err != nil {
		sendFunc(fmt.Sprintf("parse url error: %s", err.Error()))
		return "", fmt.Errorf("parse url error: %s", err.Error())
	}

	targetFile := filepath.Join(tempDir, filepath.Base(urlObject.Path))
	sendFunc(fmt.Sprintf("download file %s to %s", input.FileURL, targetFile))
	if _, err := util.DownloadToDest(input.FileURL, targetFile); err != nil {
		sendFunc(fmt.Sprintf("download error: %s", err.Error()))
		return "", fmt.Errorf("download error: %s", err.Error())
	}
	sendFunc(fmt.Sprintf("download file %s to %s success", input.FileURL, targetFile))
	remoteFile := filepath.Join(tempRemoteDir, filepath.Base(urlObject.Path))

	// 2. upload file
	sendFunc(fmt.Sprintf("upload file %s to %s....", targetFile, remoteFile))
	if err := sshClient.UploadTo(targetFile, remoteFile); err != nil {
		sendFunc(fmt.Sprintf("upload error: %s", err.Error()))
		return "", fmt.Errorf("upload error: %s", err.Error())
	}
	sendFunc(fmt.Sprintf("upload file %s to %s success", targetFile, remoteFile))

	// 3. import file
	sendFunc(fmt.Sprintf("import file %s to %s....", remoteFile))
	if err := sshClient.K3sImport(remoteFile); err != nil {
		sendFunc(fmt.Sprintf("import error: %s", err.Error()))
		return "", fmt.Errorf("import error: %s", err.Error())
	}
	sendFunc(fmt.Sprintf("import file %s to %s success", remoteFile, input.Deployment))

	// 4. rollout deployment
	if len(input.Namespace) > 0 && len(input.Deployment) > 0 {
		sendFunc(fmt.Sprintf("rollout deployment %s in namespace %s....", input.Deployment, input.Namespace))
		command := fmt.Sprintf("kubectl rollout restart deployment %s -n %s", input.Deployment, input.Namespace)
		if err := sshClient.RunCommand(command); err != nil {
			sendFunc(fmt.Sprintf("rollout error: %s", err.Error()))
			return "", fmt.Errorf("rollout error: %s", err.Error())
		}
		sendFunc(fmt.Sprintf("rollout deployment %s in namespace %s success", input.Deployment, input.Namespace))
	}
	return remoteFile, nil
}
