package main

import (
	"crypto/platform/app"
	"crypto/platform/collectors"
	"crypto/platform/db"
	"crypto/platform/telegram"
	"crypto/platform/utils"
	"sync"
)

func main() {
	logger := utils.NewLogger()
	db := db.NewInMemoryDBWithIDGen()
	a := app.NewApp(logger, db)

	run(a)
}

func run(a *app.App) {
	var wg sync.WaitGroup

	for _, f := range getToRun() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f(a)
		}()
	}

	wg.Wait()
}

func getToRun() [2]func(*app.App) {
	return [2]func(*app.App){
		startCollecting,
		startTgBot,
	}
}

func startCollecting(a *app.App) {
	c := collectors.NewRateCollector(a)
	toRun := func() error { return c.Run() }
	if err := utils.LogProcess(*a.Logger, "collecting", toRun); err != nil {
		a.Logger.Error("failed collect logs", "error", err)
	}
}

func startTgBot(a *app.App) {
	b := telegram.ModePooling.NewBot(a)
	toRun := func() error { return b.Run() }
	if err := utils.LogProcess(*a.Logger, "tg bot", toRun); err != nil {
		a.Logger.Error("failed run Telegram Bot", "error", err)
	}
}
