package utils

import (
	"log/slog"
	"os"
)

func InitLogger() (*os.File, error) {
	logFile, err := OpenLogFile()
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(logFile, nil))
	slog.SetDefault(logger)
	return logFile, nil
}
