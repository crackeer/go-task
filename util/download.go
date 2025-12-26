package util

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DownloadTo(urlString string, targetDir string, keepPath bool) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 忽略证书验证
			},
		},
	}

	urlString = strings.TrimSpace(urlString)
	if strings.HasPrefix(urlString, "//") {
		urlString = "https:" + urlString
	}
	object, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	target := filepath.Join(targetDir, object.Path)
	if !keepPath {
		target = filepath.Join(targetDir, filepath.Base(object.Path))
	}
	if info, err := os.Stat(target); err == nil && info.Size() > 0 {
		return target, nil
	}
	dir, _ := filepath.Split(target)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	response, err := client.Get(strings.TrimSpace(urlString))
	if err != nil {
		return "", fmt.Errorf("download error: %s", err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download error: %s", response.Status)
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return target, os.WriteFile(target, bytes, os.ModePerm)
}
