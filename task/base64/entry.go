package base64

import (
	"encoding/json"
	"fmt"
	"go-task/util"
)

// Base64Task 实现了Tool接口的演示任务
type Base64Input struct {
	Input string `json:"input"`
	Type  string `json:"type"`
}

func Run(input string, sendFunc func(string)) (string, error) {
	var d Base64Input
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return "", err
	}

	// 如果有输入，处理输入
	sendFunc(fmt.Sprintf("输入: %s", d.Input))
	var output string
	if d.Type == "encode" {
		output = util.Base64Encode(d.Input)
		sendFunc(fmt.Sprintf("编码结果: %s", output))
	} else {
		output = util.Base64Decode(d.Input)
		sendFunc(fmt.Sprintf("解码结果: %s", output))
	}
	return output, nil
}
