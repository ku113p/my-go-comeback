package main

import (
	"context"
	"crypto/platform/coinmarketcap"
	"crypto/platform/utils"
	"log/slog"
	"sync"
	"time"
)

type Parser interface {
	Parse() error
}

type RateCollector struct {
	ctx context.Context
}

func newRateCollector() *RateCollector {
	p := RateCollector{}

	return &p
}

func (c *RateCollector) WithContext(ctx context.Context) *RateCollector {
	c.ctx = ctx
	return c
}

func (c *RateCollector) GetLogger() *slog.Logger {
	return utils.Logger(c.ctx)
}

func (c *RateCollector) Run() error {
	var wg sync.WaitGroup
	defer wg.Wait()

	done := make(chan any, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()

		getPrices(10*time.Second, *c.GetLogger(), done)
	}()

	return nil
}

func getPrices(pause time.Duration, logger slog.Logger, done <-chan any) {
	getPrices := func() {
		if prices, err := coinmarketcap.GetPrices(); err != nil {
			logger.Error("get prices", "status", "error", "error", err)
		} else {
			logger.Info("got prices", "status", "success", "prices", prices)
		}
	}

	ticker := time.NewTicker(pause)
	defer ticker.Stop()

	firstTick := make(chan any, 1)
	firstTick <- 1

	getPrices()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			getPrices()
		}
	}
}
