package downloadwork

import (
	"encoding/json"
	"fmt"
	"go-task/util"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// DemoTask 实现了Tool接口的演示任务
type Input struct {
	Work    string `json:"work"`
	SaveDir string `json:"save_dir"`
}

// Run 执行任务，通过sendFunc发送结果
func Run(input string, sendFunc func(string)) (string, error) {
	var d Input
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return "", err
	}
	// 发送任务开始消息
	sendFunc("task started")

	if err := os.MkdirAll(d.SaveDir, 0755); err != nil {
		sendFunc(fmt.Sprintf("Error creating directory: %s", err.Error()))
		return "", err
	}
	bytes := []byte(d.Work)
	baseURL := gjson.GetBytes(bytes, "base_url").String()
	if len(baseURL) < 1 {
		sendFunc("base_url is empty")
		return "", fmt.Errorf("base_url is empty")
	}

	cubeSize := getWorkCubeSize(bytes)
	sendFunc(fmt.Sprintf("cubeSize: %s", cubeSize))

	urlList := getPanoramaURLS(bytes, cubeSize)
	sendFunc(fmt.Sprintf("getPanoramaURLS Count: %d", len(urlList)))
	modelURL, newBytes := getSetModelURLS(bytes)
	urlList = append(urlList, modelURL...)
	sendFunc(fmt.Sprintf("getSetModelURLS Count: %d", len(modelURL)))

	noJsonList, jsonList := splitJson(urlList)
	// 下载tilset.json文件
	sendFunc("download tilset.json files")
	for {
		if len(jsonList) <= 0 {
			break
		}
		sendFunc(fmt.Sprintf("parse json %d files", len(jsonList)))
		data, err := parseJSONFiles(jsonList, baseURL, d.SaveDir, sendFunc)
		if err != nil {
			sendFunc(fmt.Sprintf("parse json files error: %s", err.Error()))
			return "", err
		}

		tmpNoJSON, tmpJSON := splitJson(data)
		noJsonList = append(noJsonList, tmpNoJSON...)
		jsonList = tmpJSON
	}

	for index, item := range noJsonList {
		sendFunc(fmt.Sprintf("download %s", item))
		sendFunc(fmt.Sprintf("---> download [%d/%d] %s", index+1, len(noJsonList), item))
		sendFunc(fmt.Sprintf("local path: %s", filepath.Join(d.SaveDir, item)))
		if err := downloadFile(baseURL+item, filepath.Join(d.SaveDir, item)); err != nil {
			sendFunc(fmt.Sprintf("download %s error: %s", baseURL+item, err.Error()))
			return "", err
		}
		sendFunc(fmt.Sprintf("--enditem %s", item))
	}

	sendFunc("--> write work.json")
	newBytes, err := sjson.SetBytes(newBytes, "base_url", "{BASE_URL}")
	if err != nil {
		sendFunc(fmt.Sprintf("set base_url error:%s", err.Error()))
		return "", fmt.Errorf("set base_url error:%s", err.Error())
	}

	if err := os.WriteFile(filepath.Join(d.SaveDir, "work.json"), newBytes, 0755); err != nil {
		sendFunc(fmt.Sprintf("write work.json error:%s", err.Error()))
		return "", fmt.Errorf("write work.json error:%s", err.Error())
	}
	sendFunc("--> write work.json end")
	sendFunc("--> zip work....")
	targetZip := filepath.Join(d.SaveDir, "work.zip")
	os.RemoveAll(targetZip)
	if err := util.QuickZip(d.SaveDir, targetZip); err != nil {
		sendFunc(fmt.Sprintf("zip work error:%s", err.Error()))
		return "", fmt.Errorf("zip work error:%s", err.Error())
	}

	sendFunc(fmt.Sprintf("--> zip work end, zip file: %s", targetZip))
	downloadURL := util.GenerationDownloadURL(targetZip)
	sendFunc(fmt.Sprintf("downloadURL: %s", downloadURL))
	return downloadURL, nil
}
