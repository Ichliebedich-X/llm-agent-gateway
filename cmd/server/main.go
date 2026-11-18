package main

import (
	"context"
	"errors"
	"fmt"
	"llm-agent-gateway/internal/httpapi"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"llm-agent-gateway/internal/config"
	"llm-agent-gateway/internal/logger"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Server stopped with error", "err", err)
		os.Exit(1)
	}
}

func run() error {
	//1.读取配置
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	//2.初始化 logger
	log := logger.New(cfg.LogLevel)

	//3.输出一条结构化启动日志
	log.Info("server starting",
		"service", "llm-agent-gateway",
		"env", cfg.AppEnv,
		"addr", cfg.HTTPAddr,
		"log_level", cfg.LogLevel,
		"provider", cfg.ModelProvider,
	)

	//4.创建 http Server
	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httpapi.NewRouter(),

		//HTTP超时配置
		ReadHeaderTimeout: 5 * time.Second,  //读取请求头最长时间
		ReadTimeout:       15 * time.Second, //读整个请求最长时间
		WriteTimeout:      30 * time.Second, //响应写入最长时间（以后实现 SSE 时重新评估）
		IdleTimeout:       60 * time.Second, //两次请求之间，保持空闲连接的最长时间
	}

	//5.监听退出信号
	ctx, stop := signal.NotifyContext( //收到退出信号后，这个通道会被关闭，说明该开始关闭服务了
		context.Background(),
		os.Interrupt,    //通常对应终端里的 Ctrl+C
		syscall.SIGTERM, //常见的系统终止信号
	)
	defer stop() //停止监听信号，释放相关资源

	//6.在 goroutine 中启动 HTTP 服务
	serverErr := make(chan error, 1)

	go func() {
		log.Info("http server listening", "addr", server.Addr)
		//ListenAndServe() 是阻塞调用。如果直接在主 goroutine 中执行，程序就没有机会同时等待退出信号
		serverErr <- server.ListenAndServe()
	}()

	//7.主 goroutine 通过 select 等待服务异常退出，或者收到退出信号
	select {
	case err := <-serverErr:
		// HTTP 服务自行退出
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server failed: %w", err)
		}
		return nil

	case <-ctx.Done():
		// 收到退出信号
		log.Info("shutdown signal received")
	}

	//8.有序关闭 HTTP 服务
	//设置最多等待10s的关闭等待时间，Shutdown 会停止接受新请求，并等待正在处理的请求结束；超过期限就会返回错误
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "err", err)

		//超时后尝试强制关闭
		if closeErr := server.Close(); closeErr != nil {
			log.Error("force close failed", "err", closeErr)
		}
	}

	//9.等待 HTTP服务退出
	err = <-serverErr
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server failed during shutdown:%w", err)
	}

	log.Info("server shutdown complete")
	return nil
}
