package db

import (
	"crypto/platform/models"
)

type DB interface {
	UpdatePrices([]*models.TokenPrice) error
	GetPrice(string) (*models.TokenPrice, error)
}
