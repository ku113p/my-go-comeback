package db

import (
	"crypto/platform/models"
	"fmt"

	"github.com/google/uuid"
)

type DB interface {
	UpdatePrices([]*models.TokenPrice) error
	GetPrice(string) (*models.TokenPrice, error)

	ListUsers() ([]*models.User, error)
	GetUserByID(uuid.UUID) (*models.User, error)
	GetUserByTelegramChatID(int64) (*models.User, error)
	CreateUser(*models.User) (*models.User, error)
	RemoveUser(uuid.UUID) error

	ListNotificationsBySymbol(string) ([]*models.Notification, error)
	GetNotificationByID(uuid.UUID) (*models.Notification, error)
	ListNotificationsByUserID(uuid.UUID) ([]*models.Notification, error)
	CreateNotification(*models.Notification) (*models.Notification, error)
	RemoveNotification(uuid.UUID) error
}

var ErrNotExists = fmt.Errorf("object not exists")
