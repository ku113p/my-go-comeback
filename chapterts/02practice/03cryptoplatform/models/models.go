package models

import (
	"errors"
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
	ID     *uuid.UUID
	Symbol string
	Sign   CompareSign
	Amount float64
	Text   *string
	UserID *uuid.UUID
}

func NewNotification(symbol string, sign CompareSign, amount float64, userID *uuid.UUID, text *string) *Notification {
	return &Notification{
		Symbol: symbol,
		Sign:   sign,
		Amount: amount,
		UserID: userID,
		Text:   text,
	}
}

func (n *Notification) Check(p *TokenPrice) bool {
	return n.Sign.checkFunction(n.Symbol, n.Amount)(p)
}
