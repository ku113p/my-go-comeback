package services

import (
	"crypto/platform/app"
	"crypto/platform/db"
	"crypto/platform/models"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type NotificationService struct {
	*app.App
}

func newNotificationService(a *app.App) *NotificationService {
	return &NotificationService{a}
}

type ExpectedError struct {
	Message string
}

func (e *ExpectedError) Error() string {
	return e.Message
}

func newExpectedError(msg string) error {
	return &ExpectedError{Message: msg}
}

func (s *NotificationService) Create(user *models.User, input string) (*models.Notification, error) {
	n, err := newNotificationFromString(input)
	if err != nil {
		return nil, fmt.Errorf("invalid notification format: %w", err)
	}

	token, err := s.DB.GetPrice(n.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get price: %w", err)
	}

	if n.Check(token) {
		return nil, newExpectedError("price already reached target amount")
	}

	n.UserID = user.ID
	n, err = s.DB.CreateNotification(n)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return n, nil
}

func newNotificationFromString(s string) (*models.Notification, error) {
	words := strings.SplitN(s, " ", 3)
	if len(words) != 3 {
		return nil, fmt.Errorf("invalid format")
	}

	symbol, signString, amountString := words[0], words[1], words[2]
	symbol = strings.ToUpper(symbol)

	sign, err := models.ParseSign(signString)
	if err != nil {
		return nil, fmt.Errorf("invalid sign")
	}

	amountString = strings.ReplaceAll(amountString, ",", ".")
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount")
	}

	msg := fmt.Sprintf("price %v %v %v", symbol, sign.String(), amount)
	n := models.NewNotification(symbol, *sign, amount, nil, &msg)

	return n, nil
}

func (s *NotificationService) GetByUser(user *models.User) ([]*models.Notification, error) {
	ns, err := s.DB.ListNotificationsByUserID(*user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return ns, nil
}

func (s *NotificationService) GetByID(id string) (*models.Notification, error) {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse notification id: %w", err)
	}

	n, err := s.DB.GetNotificationByID(notificationID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			return nil, newExpectedError("notification not found")
		}
		return nil, fmt.Errorf("failed to get notification by id: %w", err)
	}

	return n, nil
}

func (s *NotificationService) DeleteByID(id string) error {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse notification id: %w", err)
	}

	err = s.DB.RemoveNotification(notificationID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			return newExpectedError("notification not found")
		}

		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}
