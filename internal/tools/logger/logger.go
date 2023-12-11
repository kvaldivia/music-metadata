package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func GetLogger() *slog.Logger {
	return logger
}
