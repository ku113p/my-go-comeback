package db

import (
	"crypto/platform/models"

	"github.com/google/uuid"
)

type DB interface {
	UpdatePrices([]*models.TokenPrice) error
	GetPrice(string) (*models.TokenPrice, error)

	GetUserByID(uuid.UUID) (*models.User, error)
	GetUserByTelegramChatID(int64) (*models.User, error)
	CreateUser(*models.User) (*models.User, error)
	RemoveUser(uuid.UUID) error

	ListNotificationsBySymbol(string) ([]*models.Notification, error)
	ListNotificationsByUserID(uuid.UUID) ([]*models.Notification, error)
	CreateNotification(*models.Notification) (*models.Notification, error)
	RemoveNotification(uuid.UUID) error
}
