package main

import (
	"llm-agent-gateway/internal/httpapi"
	"log/slog"
	"net/http"
	"os"
	"time"

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

	// 4.创建 http Server
	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httpapi.NewRouter(),

		//HTTP超时配置
		ReadHeaderTimeout: 5 * time.Second,  //读取请求头最长时间
		ReadTimeout:       15 * time.Second, //读整个请求最长时间
		WriteTimeout:      30 * time.Second, //响应写入最长时间（以后实现 SSE 时重新评估）
		IdleTimeout:       60 * time.Second, //两次请求之间，保持空闲连接的最长时间
	}

	//5.启动 HTTP服务
	log.Info("http server listening", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("http server failed", "err", err)
		os.Exit(1)
	}
}
