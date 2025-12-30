package main

import (
	"fmt"
	base64Task "go-task/task/base64"
	"go-task/task/downloadfromjson"
	json2csvTask "go-task/task/json2csv"
	md5Task "go-task/task/md5"
	qrcodeTask "go-task/task/qrcode"
	s3temptokenTask "go-task/task/s3temptoken"
	sshk3simport "go-task/task/sshk3simport"
	"go-task/task/sshrun"
	"go-task/task/sshupload"
	"os"

	"github.com/crackeer/task-facade/server"
)

func main() {
	toolMapping := make(map[string]func(string, func(string)) (string, error))
	toolMapping["base64"] = base64Task.Run
	toolMapping["json2csv"] = json2csvTask.Run
	toolMapping["md5"] = md5Task.Run
	toolMapping["qrcode"] = qrcodeTask.Run
	toolMapping["downloadfromjson"] = downloadfromjson.Run
	toolMapping["sshupload"] = sshupload.Run
	toolMapping["sshk3simport"] = sshk3simport.Run
	toolMapping["s3temptoken"] = s3temptokenTask.Run
	toolMapping["sshrun"] = sshrun.Run

	args := os.Args[1:]
	// run web server
	if len(args) == 0 {
		server.Run(toolMapping, "")
		return
	}
	// run task
	tool := args[0]
	if toolMapping[tool] == nil {
		fmt.Println("Usage: ./go-task <tool> <input_file>")
		return
	}

	input := "{}"
	if len(args) > 1 {
		bytes, _ := os.ReadFile(args[1])
		input = string(bytes)
	}

	fmt.Println("tool:", tool)
	fmt.Println("input:", input)
	toolMapping[tool](input, func(msg string) {
		fmt.Println(msg)
	})
}
