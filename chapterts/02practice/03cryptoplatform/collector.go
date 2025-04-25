package main

import (
	"context"
	"crypto/platform/coinmarketcap"
	"crypto/platform/db"
	"crypto/platform/utils"
	"sync"
	"time"
)

type RateCollector struct{}

func newRateCollector() *RateCollector {
	p := RateCollector{}

	return &p
}

func (c *RateCollector) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()

		getPrices(ctx, 10*time.Second)
	}()

	return nil
}

func getPrices(ctx context.Context, pause time.Duration) {
	logger := utils.Logger(ctx)
	database := db.Database(ctx)

	getPrices := func() {
		logger.Info("get prices inited")
		prices, err := coinmarketcap.GetPrices()
		if err != nil {
			logger.Error("get prices", "error", err)
			return
		}

		database.UpdatePrices(prices)
		logger.Info("update prices")
	}

	ticker := time.NewTicker(pause)
	defer ticker.Stop()

	for {
		getPrices()
		<-ticker.C
	}
}
