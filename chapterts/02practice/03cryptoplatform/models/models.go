package models

import "time"

type TokenPrice struct {
	Price   float64
	Name    string
	Symbold string
	Time    time.Time
}

func NewTokenPrice(p float64, n, s string, t time.Time) *TokenPrice {
	return &TokenPrice{p, n, s, t}
}
