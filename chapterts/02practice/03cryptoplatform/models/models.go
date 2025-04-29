package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenPrice struct {
	Price  float64
	Name   string
	Symbol string
	Time   time.Time
}

func NewTokenPrice(p float64, n, s string, t time.Time) *TokenPrice {
	return &TokenPrice{p, n, s, t}
}

type User struct {
	ID             *uuid.UUID
	TelegramChatID *int64
}

func NewUser(id int64) *User {
	return &User{TelegramChatID: &id}
}

type Notification struct {
	ID     *uuid.UUID
	Symbol string
	Check  func(p *TokenPrice) bool
	Text   *string
	UserID *uuid.UUID
}

func NewNotification(symbol string, check func(p *TokenPrice) bool, userID *uuid.UUID, text *string) *Notification {
	return &Notification{
		Symbol: symbol,
		Check:  check,
		UserID: userID,
		Text:   text,
	}
}
