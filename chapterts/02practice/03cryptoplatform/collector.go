package main

import (
	"context"
	"crypto/platform/coinmarketcap"
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

	getPrices := func() {
		if prices, err := coinmarketcap.GetPrices(); err != nil {
			logger.Error("get prices", "status", "error", "error", err)
		} else {
			logger.Info("got prices", "status", "success", "prices", prices)
		}
	}

	ticker := time.NewTicker(pause)
	defer ticker.Stop()

	for {
		getPrices()
		<-ticker.C
	}
}
