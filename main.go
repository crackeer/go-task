package main

import (
	base64Task "go-task/task/base64"
	json2csvTask "go-task/task/json2csv"
	md5Task "go-task/task/md5"
	qrcodeTask "go-task/task/qrcode"

	"github.com/crackeer/task-facade/server"
)

func main() {
	toolMapping := make(map[string]func(string, func(string)) (string, error))
	toolMapping["base64"] = base64Task.Run
	toolMapping["json2csv"] = json2csvTask.Run
	toolMapping["md5"] = md5Task.Run
	toolMapping["qrcode"] = qrcodeTask.Run
	server.Run(toolMapping, "")
}
