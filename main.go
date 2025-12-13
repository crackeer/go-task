package main

import (
	"encoding/json"
	"fmt"
	"go-task/container"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	base64Task "go-task/task/base64"
	downloadworkTask "go-task/task/downloadwork"
	json2csvTask "go-task/task/json2csv"
	md5Task "go-task/task/md5"
	qrcodeTask "go-task/task/qrcode"
	"go-task/util"

	_ "github.com/joho/godotenv/autoload"
)

var (
	toolMapping  map[string]func(string, func(string)) (string, error) = make(map[string]func(string, func(string)) (string, error))
	inputMapping map[string]string                                     = make(map[string]string)
	locker       sync.Mutex                                            = sync.Mutex{}
)

func enableCORS(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}

func register() {
	toolMapping["base64"] = base64Task.Run
	toolMapping["downloadwork"] = downloadworkTask.Run
	toolMapping["json2csv"] = json2csvTask.Run
	toolMapping["md5"] = md5Task.Run
	toolMapping["qrcode"] = qrcodeTask.Run
}

func main() {
	// 初始化配置
	if err := container.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	register()
	cfg := container.GetConfig()
	// 创建 Gin 实例
	router := gin.Default()
	router.Use(enableCORS)
	router.POST("/input", createInput)
	router.GET("/run/:tool", runTask)

	// 启动服务器
	router.Run(":" + cfg.Port)
}

func createInput(ctx *gin.Context) {
	inputKey := fmt.Sprintf("%d", time.Now().UnixNano())
	bytes, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to get raw data: %v", err)})
		return
	}
	locker.Lock()
	defer locker.Unlock()
	inputMapping[inputKey] = string(bytes)
	ctx.JSON(http.StatusOK, gin.H{"input_key": inputKey})
}

func runTask(ctx *gin.Context) {

	// 设置SSE头部
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")
	tool := ctx.Param("tool")

	var (
		result string
		err    error
	)
	defer func() {
		closeSSE(ctx, result)
	}()

	var printMessage func(string) = newPrintMessage(ctx)

	if tool == "" {
		printMessage("tool is required")
		return
	}

	toolFunc, ok := toolMapping[tool]
	if !ok {
		printMessage(fmt.Sprintf("tool %s not found", tool))
		return
	}

	input := getInput(ctx)

	result, err = toolFunc(input, printMessage)
	if err != nil {
		printMessage("")
		printMessage(fmt.Sprintf("failed to run task: %v", err))
		return
	}
	if len(result) > 0 {
		printMessage("output: " + result)
	}
	printMessage("")
	printMessage("Task completed successfully")
}

func getInput(ctx *gin.Context) string {
	query := util.QueryAll(ctx)
	bytes, _ := json.Marshal(query)
	inputKey, ok := query["input_key"]
	if !ok {
		return string(bytes)
	}
	if value, ok := inputMapping[inputKey]; ok {
		return value
	}

	return string(bytes)
}

func newPrintMessage(ctx *gin.Context) func(string) {
	return func(msg string) {
		ctx.SSEvent("message", msg)
		ctx.Writer.Flush()
	}
}

func closeSSE(ctx *gin.Context, msg string) {
	ctx.SSEvent("close", msg)
	ctx.Writer.Flush()
}
