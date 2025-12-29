package sshk3simport

import (
	"encoding/json"
	"fmt"
	"go-task/util"
	"os"
	"time"
)

type K3sImportInput struct {
	Config     []string `json:"config"`
	FileURL    string   `json:"file_url"`
	Namespace  string   `json:"namespace"`
	Deployment string   `json:"deployment"`
}

// RunK3sImport
//
//	@param data
//	@param sendFunc
//	@return string
//	@return error
func Run(data string, sendFunc func(string)) (string, error) {
	var input K3sImportInput
	if err := json.Unmarshal([]byte(data), &input); err != nil {
		return "", fmt.Errorf("unmarshal input error: %s", err.Error())
	}
	tempFiles := util.DefaultSplitString(input.FileURL)
	if len(tempFiles) == 0 {
		sendFunc(fmt.Sprintf("file url error: %s", input.FileURL))
		return "", fmt.Errorf("file url error: %s", input.FileURL)
	}

	tempDir, err := os.MkdirTemp("", "task-sshimport-")
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
	remoteDir := fmt.Sprintf("/tmp/k3s-image-%s-%d", input.Deployment, time.Now().Unix())
	for _, item := range input.Config {
		sendFunc(fmt.Sprintf("new ssh client %s", item))
		sshClient, err := util.NewSSHClientByConfig(item)
		if err != nil {
			sendFunc(fmt.Sprintf("new ssh client error: %s", err.Error()))
			continue
		}
		sendFunc(fmt.Sprintf("upload to %s", item))
		remoteFiles, err := util.SSHUpload(sshClient, localFiles, remoteDir, sendFunc)
		if err != nil {
			sendFunc(fmt.Sprintf("upload error: %s", err.Error()))
			continue
		}
		for _, remoteFile := range remoteFiles {
			sendFunc(fmt.Sprintf("import image %s to %s", remoteFile, item))
			output, err := sshClient.K3sImport(remoteFile)
			if err != nil {
				sendFunc(fmt.Sprintf("import error: %s", err.Error()))
				continue
			}
			sendFunc(fmt.Sprintf("import output: %s", output))
		}
		runCmd := fmt.Sprintf("kubectl rollout restart deployment %s -n %s", input.Deployment, input.Namespace)
		if len(input.Deployment) == 0 {
			runCmd = fmt.Sprintf("kubectl rollout restart deployment -n %s", input.Namespace)
		}
		sendFunc(fmt.Sprintf("run command %s", runCmd))
		output, err := sshClient.Run(runCmd)
		if err != nil {
			sendFunc(fmt.Sprintf("run command error: %s", err.Error()))
			continue
		}
		sendFunc(fmt.Sprintf("run command output: %s", output))

		sshClient.Close()
		retData = append(retData, remoteFiles)
	}

	byteData, err := json.Marshal(retData)
	if err != nil {
		sendFunc(fmt.Sprintf("marshal error: %s", err.Error()))
		return "", fmt.Errorf("marshal error: %s", err.Error())
	}

	return string(byteData), nil
}
