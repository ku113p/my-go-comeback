package collectors

import (
	"crypto/platform/app"
	"crypto/platform/coinmarketcap"
	"sync"
	"time"
)

const updateSeconds = 60

type RateCollector struct {
	app *app.App
}

func NewRateCollector(app *app.App) *RateCollector {
	p := RateCollector{app}

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

		c.app.DB.UpdatePrices(prices)
		c.app.Logger.Info("prices updated")
		c.app.Logger.Info("get prices finished")
	}

	ticker := time.NewTicker(pause)
	defer ticker.Stop()

	for {
		getPrices()
		<-ticker.C
	}
}
