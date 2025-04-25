package main

import (
	"crypto/platform/db"
	"crypto/platform/utils"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	logger := utils.NewLogger()
	database := db.NewInMemoryDB()

	ctx := utils.NewContext()
	ctx = utils.WithLogger(ctx, logger)
	ctx = db.WithDatabase(ctx, database)

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
