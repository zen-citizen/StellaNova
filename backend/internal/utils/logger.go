package utils

import (
	"log/slog"
	"os"
	"strings"
)

func SetupLogger(level, format string) *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	var slogHandler slog.Handler
	if strings.ToLower(format) == "json" {
		slogHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		slogHandler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(slogHandler)
}
