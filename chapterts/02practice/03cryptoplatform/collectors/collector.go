package collectors

import (
	"crypto/platform/app"
	"crypto/platform/coinmarketcap"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

const updateSeconds = 60

type RateCollector struct {
	app     *app.App
	updated chan<- struct{}
}

func NewRateCollector(app *app.App, updated chan<- struct{}) *RateCollector {
	p := RateCollector{app, updated}

	return &p
}

func (c *RateCollector) Run() error {
	var wg sync.WaitGroup
	defer wg.Wait()

	interval, err := updateInteraval()
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		c.getPrices(*interval)
	}()

	return nil
}

func (c *RateCollector) getPrices(pause time.Duration) {
	getPrices := func() {
		c.app.Logger.Info("get prices inited")

		prices, err := coinmarketcap.GetPrices()
		if err != nil {
			c.app.Logger.Error("get prices", "error", err)
			return
		}

		if err := c.app.DB.UpdatePrices(prices); err != nil {
			c.app.Logger.Error("failed update prices", "error", err)
			return
		}

		c.app.Logger.Info("prices updated")
		c.app.Logger.Info("get prices finished")
		go c.notifiUpdated()
	}

	ticker := time.NewTicker(pause)
	defer ticker.Stop()

	for {
		getPrices()
		<-ticker.C
	}
}

func (c *RateCollector) notifiUpdated() {
	c.updated <- struct{}{}
}

func updateInteraval() (*time.Duration, error) {
	intervalStr, ok := os.LookupEnv("UPDATE_INTERVAL_MS")
	if !ok {
		return nil, fmt.Errorf("environment variable UPDATE_INTERVAL_MS not set")
	}

	intervalMs, err := strconv.Atoi(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid UPDATE_INTERVAL_MS value: %v", err)
	}

	interval := time.Duration(intervalMs) * time.Millisecond

	return &interval, nil
}
