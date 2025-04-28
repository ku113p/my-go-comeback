package utils

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}

func LogProcess[T any](toCall func() T, name string, l slog.Logger) T {
	l.Info(name, "status", "started")
	result := toCall()
	l.Info(name, "status", "finished")

	return result
}
