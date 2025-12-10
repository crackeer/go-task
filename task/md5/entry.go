package md5

import (
	"encoding/json"
	"fmt"
	"go-task/util"
)

// Md5Task 实现了Tool接口的演示任务
type Md5Task struct {
	Input string `json:"input"`
}

// Run 执行任务，通过sendFunc发送结果
func Run(input string, sendFunc func(string)) (string, error) {
	var d Md5Task
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return "", err
	}
	// 如果有输入，处理输入
	sendFunc(fmt.Sprintf("输入: %s", d.Input))
	// 计算MD5值
	md5Value := util.CalculateMD5(d.Input)
	sendFunc(fmt.Sprintf("MD5值: %s", md5Value))
	return md5Value, nil
}
