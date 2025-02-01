package utils

import (
	"fmt"
	"log/slog"
)

func InitLogger() {
	logFile, err := OpenLogFile()
	if err != nil {
		slog.Error("Could not create or open log file, defaulting to standard error")
		slog.Error(fmt.Sprintf("%v\n", err))
	} else {
		logger := slog.New(slog.NewJSONHandler(logFile, nil))
		slog.SetDefault(logger)
	}
}
