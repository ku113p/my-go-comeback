package collectors

import (
	"crypto/platform/app"
	"crypto/platform/coinmarketcap"
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

	wg.Add(1)
	go func() {
		defer wg.Done()

		c.getPrices(updateSeconds * time.Second)
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
