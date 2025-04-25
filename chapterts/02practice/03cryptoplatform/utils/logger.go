package utils

import (
	"context"
	"log/slog"
	"os"
)

type ctxLogger string

const loggerKey ctxLogger = "logger"

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func Logger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}

type RunnableHasLogger[T any] interface {
	GetLogger() *slog.Logger
	Run() T
}

func WrapRunningLog[T any](hl RunnableHasLogger[T], name string) T {
	l := hl.GetLogger()
	l.Info(name, "status", "started")
	result := hl.Run()
	l.Info(name, "status", "finished")

	return result
}
