package qrcode

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/skip2/go-qrcode"
)

// QrcodeTask 实现了Tool接口的演示任务
type QrcodeTask struct {
	Input string `json:"input"`
}

// Run 执行任务，通过sendFunc发送结果
func Run(input string, sendFunc func(string)) (string, error) {
	var d QrcodeTask
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return "", err
	}
	// 如果有输入，处理输入
	sendFunc(fmt.Sprintf("输入: %s", d.Input))
	qr, err := qrcode.New(d.Input, qrcode.High)
	if err != nil {
		return "", err
	}
	qrASCII := qr.ToSmallString(false)
	parts := strings.Split(qrASCII, "\n")
	for _, part := range parts {
		sendFunc(part)
	}
	return qrASCII, nil
}
