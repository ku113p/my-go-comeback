package db

import (
	"context"
	"crypto/platform/models"
)

type DB interface {
	UpdatePrices([]*models.TokenPrice) error
	GetPrice(string) (*models.TokenPrice, error)
}

type ctxDatabase string

const databaseKey ctxDatabase = "database"

func WithDatabase(ctx context.Context, database DB) context.Context {
	return context.WithValue(ctx, databaseKey, database)
}

func Database(ctx context.Context) DB {
	if logger, ok := ctx.Value(databaseKey).(DB); ok {
		return logger
	}
	return nil
}
