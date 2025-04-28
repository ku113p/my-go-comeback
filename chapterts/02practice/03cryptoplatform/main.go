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

	l := utils.NewLogger()
	d := db.NewInMemoryDB()

	a := app.NewApp(l, d)

	c := collectors.NewRateCollector(a)
	wg.Add(1)
	go func() {
		defer wg.Done()

		toRun := func() error { return c.Run() }
		if err := app.LogProcess(a, "collecting", toRun); err != nil {
			l.Error("failed collect logs", "error", err)
		}
	}()

	wg.Wait()
}
