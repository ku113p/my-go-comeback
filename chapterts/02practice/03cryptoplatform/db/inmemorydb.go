package db

import (
	"crypto/platform/models"
	"fmt"

	"github.com/google/uuid"
)

type InMemoryDB struct {
	tokensStorage        map[string]*models.TokenPrice
	usersStorage         map[uuid.UUID]*models.User
	notificationsStorage map[uuid.UUID]*models.Notification

	idGenerator func() uuid.UUID
	locker      chan any
}

func NewInMemoryDBWithIDGen() *InMemoryDB {
	genID := func() uuid.UUID {
		return uuid.New()
	}

	return newInMemoryDB(genID)
}

func newInMemoryDB(idGenerator func() uuid.UUID) *InMemoryDB {
	return &InMemoryDB{
		tokensStorage:        make(map[string]*models.TokenPrice, 0),
		usersStorage:         make(map[uuid.UUID]*models.User, 0),
		notificationsStorage: make(map[uuid.UUID]*models.Notification, 0),

		idGenerator: idGenerator,
		locker:      make(chan any, 1),
	}
}

func (db *InMemoryDB) UpdatePrices(newPirces []*models.TokenPrice) error {
	db.locker <- nil
	defer func() { <-db.locker }()

	newStorage := make(map[string]*models.TokenPrice, len(newPirces))
	for _, p := range newPirces {
		newStorage[p.Symbol] = p
	}
	db.tokensStorage = newStorage
	return nil
}

func (db *InMemoryDB) GetPrice(symbol string) (*models.TokenPrice, error) {
	db.locker <- nil
	defer func() { <-db.locker }()

	tp, ok := db.tokensStorage[symbol]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	return tp, nil
}

func (db *InMemoryDB) CreateNotification(n *models.Notification) (*models.Notification, error) {
	db.locker <- nil
	defer func() { <-db.locker }()

	n.ID = db.newID()
	db.notificationsStorage[*n.ID] = n

	return n, nil
}

func (db *InMemoryDB) newID() *uuid.UUID {
	id := db.idGenerator()
	return &id
}

func (db *InMemoryDB) CreateUser(u *models.User) (*models.User, error) {
	db.locker <- nil
	defer func() { <-db.locker }()

	u.ID = db.newID()
	db.usersStorage[*u.ID] = u

	return u, nil
}

func (db *InMemoryDB) GetUserByID(id uuid.UUID) (*models.User, error) {
	db.locker <- nil
	defer func() { <-db.locker }()

	u, ok := db.usersStorage[id]
	if !ok {
		return nil, ErrNotExists
	}

	return u, nil
}

func (db *InMemoryDB) GetUserByTelegramChatID(id int64) (*models.User, error) {
	db.locker <- nil
	defer func() { <-db.locker }()

	for _, u := range db.usersStorage {
		if *u.TelegramChatID == id {
			return u, nil
		}
	}

	return nil, ErrNotExists
}

func (db *InMemoryDB) ListNotificationsBySymbol(symbol string) ([]*models.Notification, error) {
	suiteFunc := func(n *models.Notification) bool {
		return n.Symbol == symbol
	}
	return db.collectNotifications(suiteFunc)
}

func (db *InMemoryDB) collectNotifications(suite func(*models.Notification) bool) ([]*models.Notification, error) {
	db.locker <- nil
	defer func() { <-db.locker }()

	notifications := make([]*models.Notification, 0)

	for _, n := range db.notificationsStorage {
		if suite(n) {
			notifications = append(notifications, n)
		}
	}

	return notifications, nil
}

func (db *InMemoryDB) ListNotificationsByUserID(userID uuid.UUID) ([]*models.Notification, error) {
	suiteFunc := func(n *models.Notification) bool {
		return *n.UserID == userID
	}
	return db.collectNotifications(suiteFunc)
}

func (db *InMemoryDB) RemoveNotification(id uuid.UUID) error {
	db.locker <- nil
	defer func() { <-db.locker }()

	_, ok := db.notificationsStorage[id]
	if !ok {
		return ErrNotExists
	}

	delete(db.notificationsStorage, id)
	return nil
}

func (db *InMemoryDB) RemoveUser(id uuid.UUID) error {
	db.locker <- nil
	defer func() { <-db.locker }()

	_, ok := db.usersStorage[id]
	if !ok {
		return ErrNotExists
	}

	delete(db.usersStorage, id)
	return nil
}
