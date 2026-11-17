package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New 根据配置的日志级别创建统一的 slog.Logger。
// 使用标准库 log/slog
func New(level string) *slog.Logger {
	var lv slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lv = slog.LevelDebug
	case "warn":
		lv = slog.LevelWarn
	case "error":
		lv = slog.LevelError
	default:
		lv = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lv,
	})
	return slog.New(handler)
}
