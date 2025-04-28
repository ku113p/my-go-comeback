package db

import (
	"crypto/platform/models"
	"fmt"
)

type InMemoryDB struct {
	storage map[string]*models.TokenPrice
	locker  chan any
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		storage: make(map[string]*models.TokenPrice, 0),
		locker:  make(chan any, 1),
	}
}

func (db *InMemoryDB) UpdatePrices(newPirces []*models.TokenPrice) error {
	db.locker <- nil
	defer func() { <-db.locker }()

	newStorage := make(map[string]*models.TokenPrice, len(newPirces))
	for _, p := range newPirces {
		newStorage[p.Symbol] = p
	}
	db.storage = newStorage
	return nil
}

func (db *InMemoryDB) GetPrice(symbol string) (*models.TokenPrice, error) {
	db.locker <- nil
	defer func() { <-db.locker }()

	tp, ok := db.storage[symbol]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	return tp, nil
}
