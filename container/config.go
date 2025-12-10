package container

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type AppConfig struct {
	Port string `env:"PORT" envDefault:"80"`
}

var (
	cfg *AppConfig
)

func InitConfig() error {
	// 从环境变量中解析配置
	cfg = &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
	return nil
}

func GetConfig() *AppConfig {
	return cfg
}
