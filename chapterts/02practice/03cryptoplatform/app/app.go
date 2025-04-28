package app

import (
	"crypto/platform/db"
	"log/slog"
)

type App struct {
	Logger *slog.Logger
	DB     db.DB
}

func NewApp(logger *slog.Logger, db db.DB) *App {
	return &App{logger, db}
}
