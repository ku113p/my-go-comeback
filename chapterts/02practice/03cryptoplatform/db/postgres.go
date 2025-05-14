package db

import (
	"context"
	"crypto/platform/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type RealDB struct {
	idGenerator func() uuid.UUID
	db          sqlx.ExtContext
	Close       func()
	migrate     func() error
}

func (p *RealDB) newID() *uuid.UUID {
	id := p.idGenerator()
	return &id
}

func (p *RealDB) Migrate() error {
	return p.migrate()
}

func NewPostgresDBWithIDGen(dbURI string) (*RealDB, error) {
	genID := func() uuid.UUID {
		return uuid.New()
	}

	db, err := sqlx.Connect("postgres", dbURI)
	if err != nil {
		return nil, err
	}
	closeFunc := func() { db.Close() }
	migrateFunc := func() error { return migratePostgreSQL(db.DB) }

	return &RealDB{genID, db, closeFunc, migrateFunc}, nil
}

func (p *RealDB) UpdatePrices(prices []*models.TokenPrice) error {
	_, err := sqlx.NamedExecContext(context.TODO(), p.db, `
		INSERT INTO token_price (price, name, symbol, time)
		VALUES (:price, :name, :symbol, :time)
		ON CONFLICT (symbol) DO UPDATE
		SET price = EXCLUDED.price,
			name = EXCLUDED.name,
			time = EXCLUDED.time;`, prices)
	return err
}

func (p *RealDB) GetPrice(symbol string) (*models.TokenPrice, error) {
	prices := []*models.TokenPrice{}
	err := sqlx.SelectContext(context.TODO(), p.db, &prices, `SELECT * FROM token_price WHERE symbol = $1`, symbol)
	if err != nil {
		return nil, err
	}
	switch len(prices) {
	case 0:
		return nil, ErrNotExists
	case 1:
		return prices[0], nil
	}
	return nil, fmt.Errorf("too much prices")
}

func (p *RealDB) ListUsers() ([]*models.User, error) {
	var users []*models.User
	err := sqlx.SelectContext(context.TODO(), p.db, &users, `SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (p *RealDB) GetUserByID(id uuid.UUID) (*models.User, error) {
	users := []*models.User{}
	err := sqlx.SelectContext(context.TODO(), p.db, &users, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	switch len(users) {
	case 0:
		return nil, ErrNotExists
	case 1:
		return users[0], nil
	}
	return nil, fmt.Errorf("too much users")
}

func (p *RealDB) GetUserByTelegramChatID(telegramChatID int64) (*models.User, error) {
	users := []*models.User{}
	err := sqlx.SelectContext(context.TODO(), p.db, &users, `SELECT * FROM users WHERE telegram_chat_id = $1`, telegramChatID)
	if err != nil {
		return nil, err
	}
	switch len(users) {
	case 0:
		return nil, ErrNotExists
	case 1:
		return users[0], nil
	}
	return nil, fmt.Errorf("too much users")
}

func (p *RealDB) CreateUser(user *models.User) (*models.User, error) {
	user.ID = p.newID()
	_, err := sqlx.NamedExecContext(context.TODO(), p.db, `INSERT INTO users (id, telegram_chat_id) VALUES (:id, :telegram_chat_id)`, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *RealDB) RemoveUser(id uuid.UUID) error {
	_, err := p.db.ExecContext(context.TODO(), `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (p *RealDB) ListNotificationsBySymbol(symbol string) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := sqlx.SelectContext(context.TODO(), p.db, &notifications, `SELECT * FROM notification WHERE symbol = $1`, symbol)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (p *RealDB) GetNotificationByID(id uuid.UUID) (*models.Notification, error) {
	notifications := []*models.Notification{}
	err := sqlx.SelectContext(context.TODO(), p.db, &notifications, `SELECT * FROM notification WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	switch len(notifications) {
	case 0:
		return nil, ErrNotExists
	case 1:
		return notifications[0], nil
	}
	return nil, fmt.Errorf("too much notifications")
}

func (p *RealDB) ListNotificationsByUserID(id uuid.UUID) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := sqlx.SelectContext(context.TODO(), p.db, &notifications, `SELECT * FROM notification WHERE user_id = $1`, id)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (p *RealDB) CreateNotification(n *models.Notification) (*models.Notification, error) {
	n.ID = p.newID()
	_, err := sqlx.NamedExecContext(context.TODO(), p.db, `INSERT INTO notification (id, user_id, symbol, sign, amount) VALUES (:id, :user_id, :symbol, :sign, :amount)`, n)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (p *RealDB) RemoveNotification(id uuid.UUID) error {
	_, err := p.db.ExecContext(context.TODO(), `DELETE FROM notification WHERE id = $1`, id)
	return err
}
