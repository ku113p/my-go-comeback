package main

import (
	"crypto/platform/app"
	"crypto/platform/collectors"
	"crypto/platform/db"
	"crypto/platform/utils"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	logger := utils.NewLogger()
	db := db.NewInMemoryDBWithIDGen()

	a := app.NewApp(logger, db)

	c := collectors.NewRateCollector(a)
	wg.Add(1)
	go func() {
		defer wg.Done()

		toRun := func() error { return c.Run() }
		if err := app.LogProcess(a, "collecting", toRun); err != nil {
			logger.Error("failed collect logs", "error", err)
		}
	}()

	wg.Wait()
}
