package db

import (
	"crypto/platform/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	idGenerator func() uuid.UUID
	db          *sqlx.DB
}

func (p *PostgresDB) newID() *uuid.UUID {
	id := p.idGenerator()
	return &id
}

func (p *PostgresDB) Close() error {
	if err := p.db.Close(); err != nil {
		return err
	}
	return nil
}

func NewPostgresDBWithIDGen(dbURI string) (*PostgresDB, error) {
	genID := func() uuid.UUID {
		return uuid.New()
	}

	return newPostgresDB(genID, dbURI)
}

func newPostgresDB(idGenerator func() uuid.UUID, dbURI string) (*PostgresDB, error) {
	db, err := sqlx.Connect("postgres", dbURI)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{idGenerator, db}, nil
}

func (db *PostgresDB) Migrate() error {
	return migratePostgreSQL(db.db.DB)
}

func (p *PostgresDB) UpdatePrices(prices []*models.TokenPrice) error {
	_, err := p.db.NamedExec(`
	INSERT INTO token_price (price, name, symbol, time)
	VALUES (:price, :name, :symbol, :time)
	ON CONFLICT (symbol) DO UPDATE
	SET price = EXCLUDED.price,
		name = EXCLUDED.name,
		time = EXCLUDED.time;`, prices)
	return err
}

func (p *PostgresDB) GetPrice(symbol string) (*models.TokenPrice, error) {
	prices := []*models.TokenPrice{}
	err := p.db.Select(&prices, `SELECT * FROM token_price WHERE symbol = $1`, symbol)
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

func (p *PostgresDB) ListUsers() ([]*models.User, error) {
	var users []*models.User
	err := p.db.Select(&users, `SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (p *PostgresDB) GetUserByID(id uuid.UUID) (*models.User, error) {
	users := []*models.User{}
	err := p.db.Select(&users, `SELECT * FROM users WHERE id = $1`, id)
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

func (p *PostgresDB) GetUserByTelegramChatID(telegramChatID int64) (*models.User, error) {
	users := []*models.User{}
	err := p.db.Select(&users, `SELECT * FROM users WHERE telegram_chat_id = $1`, telegramChatID)
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

func (p *PostgresDB) CreateUser(user *models.User) (*models.User, error) {
	user.ID = p.newID()

	_, err := p.db.NamedExec(`INSERT INTO users (id, telegram_chat_id) VALUES (:id, :telegram_chat_id)`, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PostgresDB) RemoveUser(id uuid.UUID) error {
	_, err := p.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) ListNotificationsBySymbol(symbol string) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := p.db.Select(&notifications, `SELECT * FROM notification WHERE symbol = $1`, symbol)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (p *PostgresDB) GetNotificationByID(id uuid.UUID) (*models.Notification, error) {
	notifications := []*models.Notification{}
	err := p.db.Select(&notifications, `SELECT * FROM notification WHERE id = $1`, id)
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

func (p *PostgresDB) ListNotificationsByUserID(id uuid.UUID) ([]*models.Notification, error) {
	var notifications []*models.Notification
	err := p.db.Select(&notifications, `SELECT * FROM notification WHERE user_id = $1`, id)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (p *PostgresDB) CreateNotification(n *models.Notification) (*models.Notification, error) {
	n.ID = p.newID()

	_, err := p.db.NamedExec(`INSERT INTO notification (id, user_id, symbol, sign, amount) VALUES (:id, :user_id, :symbol, :sign, :amount)`, n)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (p *PostgresDB) RemoveNotification(id uuid.UUID) error {
	_, err := p.db.Exec(`DELETE FROM notification WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
