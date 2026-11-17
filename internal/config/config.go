package config

import (
	"fmt"
	"os"
	"strings"
)

// Config 是程序运行所需配置的集中定义。
// 所有字段都来自环境变量，便于本地、Docker、测试环境切换。
type Config struct {
	AppEnv        string // dev / test / prod
	HTTPAddr      string // 监听地址，如 :8080
	LogLevel      string // debug / info / warn / error
	ModelProvider string // mock / openai / ...
	ModelAPIKey   string // 真实密钥只从环境变量读取，不落日志
}

// Load 读取环境变量、补默认值并校验。
// 校验失败立即返回 error，避免"服务启动了但行为异常"。
func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:        getEnv("APP_ENV", "dev"),
		HTTPAddr:      getEnv("HTTP_ADDR", ":8080"),
		LogLevel:      strings.ToLower(getEnv("LOG_LEVEL", "info")),
		ModelProvider: strings.ToLower(getEnv("MODEL_PROVIDER", "mock")),
		ModelAPIKey:   os.Getenv("MODEL_API_KEY"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("invalid LOG_LEVEL: %s", c.LogLevel)
	}
	if c.HTTPAddr == "" {
		return fmt.Errorf("HTTP_ADDR must not be empty")
	}
	return nil
}

// getEnv 读取环境变量；为空时返回 fallback。
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
