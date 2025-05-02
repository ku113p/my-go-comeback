package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TokenPrice struct {
	Price  float64   `db:"price"`
	Name   string    `db:"name"`
	Symbol string    `db:"symbol"`
	Time   time.Time `db:"time"`
}

func NewTokenPrice(p float64, n, s string, t time.Time) *TokenPrice {
	return &TokenPrice{p, n, s, t}
}

type User struct {
	ID             *uuid.UUID `db:"id"`
	TelegramChatID *int64     `db:"telegram_chat_id"`
}

func NewUser(id int64) *User {
	return &User{TelegramChatID: &id}
}

type CompareSign string

const (
	moreSign CompareSign = ">"
	lessSign CompareSign = "<"
)

func ParseSign(s string) (*CompareSign, error) {
	sign := CompareSign(s)
	switch sign {
	case moreSign, lessSign:
		return &sign, nil
	default:
		return nil, errors.New("invalid sign")
	}
}

func (s *CompareSign) checkFunction(symbol string, amount float64) func(p *TokenPrice) bool {
	return func(p *TokenPrice) bool {
		if p.Symbol != symbol {
			return false
		}

		switch *s {
		case moreSign:
			return p.Price > amount
		case lessSign:
			return p.Price < amount
		}

		return false
	}
}

func (s *CompareSign) String() string {
	switch *s {
	case moreSign:
		return ">"
	case lessSign:
		return "<"
	}

	return "?"
}

func (s *CompareSign) When() string {
	switch *s {
	case moreSign:
		return "Got bigger"
	case lessSign:
		return "Got smaller"
	}

	return "?"
}

type Notification struct {
	ID     *uuid.UUID  `db:"iq"`
	Symbol string      `db:"symbol"`
	Sign   CompareSign `db:"sign"`
	Amount float64     `db:"amount"`
	UserID *uuid.UUID  `db:"user_id"`
}

func NewNotification(symbol string, sign CompareSign, amount float64, userID *uuid.UUID) *Notification {
	return &Notification{
		Symbol: symbol,
		Sign:   sign,
		Amount: amount,
		UserID: userID,
	}
}

func (n *Notification) Check(p *TokenPrice) bool {
	return n.Sign.checkFunction(n.Symbol, n.Amount)(p)
}
