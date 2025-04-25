package db

import (
	"crypto/platform/models"
	"fmt"
)

type InMemoryDB struct {
	storage []*models.TokenPrice
	locker  chan any
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		storage: make([]*models.TokenPrice, 0),
		locker:  make(chan any, 1),
	}
}

func (db *InMemoryDB) UpdatePrices(newPirces []*models.TokenPrice) error {
	db.storage = newPirces
	return nil
}

func (db *InMemoryDB) GetPrice(symbol string) (*models.TokenPrice, error) {
	for _, tp := range db.storage {
		if tp.Symbol == symbol {
			return tp, nil
		}
	}

	return nil, fmt.Errorf("not found")
}
