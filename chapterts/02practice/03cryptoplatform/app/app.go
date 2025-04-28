package app

import (
	"crypto/platform/db"
	"crypto/platform/utils"
	"log/slog"
)

type App struct {
	Logger *slog.Logger
	DB     db.DB
}

func NewApp(logger *slog.Logger, db db.DB) *App {
	return &App{logger, db}
}

func LogProcess[T any](app *App, name string, toCall func() T) T {
	return utils.LogProcess(toCall, name, *app.Logger)
}
