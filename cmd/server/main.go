package main

import (
	"log/slog"
	"os"

	"llm-agent-gateway/internal/config"
	"llm-agent-gateway/internal/logger"
)

func main() {
	// 1. 读取并校验配置；失败立即退出，不留半启动状态。
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config failed", "err", err)
		os.Exit(1)
	}

	// 2. 用配置里的级别创建 logger。
	log := logger.New(cfg.LogLevel)

	// 3. 输出一条结构化启动日志
	log.Info("server starting",
		"service", "llm-agent-gateway",
		"env", cfg.AppEnv,
		"addr", cfg.HTTPAddr,
		"log_level", cfg.LogLevel,
		"provider", cfg.ModelProvider,
	)

	// 4. 不启动 HTTP，先跑通骨架
	log.Info("d1 skeleton ready, http server not started yet")
}
