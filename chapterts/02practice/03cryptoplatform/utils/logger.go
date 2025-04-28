package utils

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}

func LogProcess[T any](l slog.Logger, name string, toCall func() T) T {
	l.Info(name, "status", "started")
	result := toCall()
	l.Info(name, "status", "finished")

	return result
}
