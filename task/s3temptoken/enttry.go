package s3temptoken

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Input struct {
	Endpoint   string `json:"endpoint"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Bucket     string `json:"bucket"`
	ObjectName string `json:"object_name"`
}

func Run(input string, sendFunc func(string)) (string, error) {
	var inputObj Input
	if err := json.Unmarshal([]byte(input), &inputObj); err != nil {
		sendFunc(fmt.Sprintf("unmarshal input error: %s", err.Error()))
		return "", err
	}
	// === 配置 MinIO STS 端点和长期凭证 ===
	minioEndpoint := inputObj.Endpoint // 替换为你的 MinIO 地址
	accessKey := inputObj.AccessKey    // 具有 AssumeRole 权限的用户
	secretKey := inputObj.SecretKey    // 对应密码
	region := "us-east-1"              // MinIO 通常忽略 region，但 SDK 需要

	// 创建静态凭证提供者
	staticCreds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	// 加载配置（自定义 endpoint）
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(staticCreds),
		config.WithBaseEndpoint(minioEndpoint),
	)
	if err != nil {
		sendFunc(fmt.Sprintf("load aws config error: %s", err.Error()))
		return "", err
	}

	// 创建 STS 客户端
	stsClient := sts.NewFromConfig(cfg)

	// 调用 AssumeRole
	arn := fmt.Sprintf("arn:aws:s3:::%s/%s", inputObj.Bucket, inputObj.ObjectName)
	durationSeconds := int32(3600) // 最大不超过 MinIO 的 max-session-duration（默认 3600 秒）
	result, err := stsClient.AssumeRole(context.TODO(), &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn), // MinIO 中可任意填写，会被忽略
		RoleSessionName: aws.String("web-upload-session"),
		DurationSeconds: &durationSeconds,
	})
	if err != nil {
		sendFunc(fmt.Sprintf("assume role error: %s", err.Error()))
		return "", err
	}

	// 输出临时凭证
	creds := result.Credentials
	sendFunc("临时凭证获取成功")
	sendFunc(fmt.Sprintf("AccessKeyId=%s, SecretAccessKey=%s, SessionToken=%s, Expiration=%s",
		*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken, creds.Expiration.Format("2006-01-02 15:04:05")))
	bytes, err := json.Marshal(creds)
	if err != nil {
		sendFunc(fmt.Sprintf("marshal creds error: %s", err.Error()))
		return "", err
	}
	return string(bytes), nil
}
