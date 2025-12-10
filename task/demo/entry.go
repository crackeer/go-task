package demo

import (
	"encoding/json"
	"fmt"
	"time"
)

type DemoInput struct {
	Input string `json:"input"`
}

// Run 执行任务，通过sendFunc发送结果
func Run(input string, sendFunc func(string)) (string, error) {
	var d DemoInput
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return "", err
	}
	// 发送任务开始消息
	sendFunc("Demo task started")

	// 如果有输入，处理输入
	if d.Input != "" {
		sendFunc(fmt.Sprintf("Processing input: %s", d.Input))
	}

	// 模拟处理步骤
	for i := 0; i < 5; i++ {
		sendFunc(fmt.Sprintf("Processing step %d/5", i+1))
		// 模拟处理延迟
		time.Sleep(time.Second)
	}

	// 发送任务完成消息
	sendFunc(fmt.Sprintf("Demo task completed successfully with input: %s", d.Input))
	return d.Input, nil
}
