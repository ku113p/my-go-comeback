package main

import (
	"crypto/platform/utils"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	logger := utils.NewLogger()
	ctx := utils.NewContext()
	ctx = utils.WithLogger(ctx, logger)

	c := newRateCollector()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := utils.LogRun(ctx, c, "collecting"); err != nil {
			logger.Error("failed collect logs", "error", err)
		}
	}()

	wg.Wait()
}
