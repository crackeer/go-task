package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type UploadFileResponse struct {
	URL string `json:"url"`
}

// UploadFile
//
//	@param filePath
//	@param destURL
//	@return UploadFileResponse
//	@return error
func UploadFile(filePath string) (UploadFileResponse, error) {
	destURL := os.Getenv("UPLOAD_URL")
	if len(destURL) == 0 {
		return UploadFileResponse{}, fmt.Errorf("upload url is empty")
	}
	resp, err := resty.New().R().SetFile("file", filePath).Post(destURL)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("upload file %s failed, err: %v", destURL, err)
	}
	if resp.StatusCode() != 200 {
		return UploadFileResponse{}, fmt.Errorf("upload file %s failed, status code: %d", filePath, resp.StatusCode())
	}

	body := resp.Body()
	response := UploadFileResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return UploadFileResponse{}, err
	}
	return response, nil
}
