package downloadurl

import (
	"encoding/json"
	"fmt"
	"go-task/util"
	"os"
	"path/filepath"
)

type Input struct {
	Data     string `json:"data"`
	KeepPath string `json:"keep_path"`
}

func Run(input string, sendFunc func(string)) (string, error) {
	var d Input
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return "", err
	}
	keepPath := d.KeepPath == "yes"

	tempDir, err := os.MkdirTemp("", "downloadfromjson-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			sendFunc(fmt.Sprintf("failed to remove temp dir: %v", err))
		}
	}()

	sendFunc(fmt.Sprintf("temp dir: %s", tempDir))
	sendFunc(fmt.Sprintf("keep path: %v", keepPath))
	sendFunc("extract url...")
	downloadURLS := util.ExtractURL(d.Data, util.IsURL)
	if len(downloadURLS) == 0 {
		sendFunc("no valid url found")
		return "", fmt.Errorf("no valid url found")
	}
	sendFunc(fmt.Sprintf("download %d urls", len(downloadURLS)))

	for _, url := range downloadURLS {
		sendFunc(fmt.Sprintf("download %s", url))
		target, err := util.DownloadTo(url, tempDir, keepPath)
		if err != nil {
			sendFunc(fmt.Sprintf("failed to download %s: %v", url, err))
			continue
		}
		sendFunc(fmt.Sprintf("download %s to %s", url, target))
	}
	zipFile := filepath.Join(tempDir, "download.zip")
	sendFunc(fmt.Sprintf("zip %s to %s", tempDir, zipFile))
	if err := util.QuickZip(tempDir, zipFile); err != nil {
		sendFunc(fmt.Sprintf("failed to zip %s: %v", tempDir, err))
		return "", fmt.Errorf("failed to zip %s: %v", tempDir, err)
	}
	sendFunc(fmt.Sprintf("zip %s success", zipFile))

	sendFunc(fmt.Sprintf("upload %s...", zipFile))
	uploadURL, err := util.UploadFile(zipFile)
	if err != nil {
		sendFunc(fmt.Sprintf("failed to upload %s: %v", zipFile, err))
		return "", fmt.Errorf("failed to upload %s: %v", zipFile, err)
	}
	sendFunc(fmt.Sprintf("upload %s to %s", zipFile, uploadURL))

	return uploadURL.URL, nil
}
