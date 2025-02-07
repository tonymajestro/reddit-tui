package utils

import (
	"log/slog"
	"os"
	"strings"
)

func InitLogger(logLevel string) (*os.File, error) {
	var level slog.Level

	logFile, err := OpenLogFile()
	if err != nil {
		return nil, err
	}

	switch l := strings.ToLower(logLevel); l {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	options := slog.HandlerOptions{Level: level}
	logger := slog.New(slog.NewJSONHandler(logFile, &options))
	slog.SetDefault(logger)
	return logFile, nil
}
