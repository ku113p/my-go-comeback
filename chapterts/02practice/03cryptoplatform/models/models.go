package models

import "time"

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
	telegramChatID int64
	notifies       map[string][]func(float64) bool
}
