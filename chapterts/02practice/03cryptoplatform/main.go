package main

import (
	"crypto/platform/app"
	"crypto/platform/collectors"
	"crypto/platform/db"
	"crypto/platform/monitoring"
	"crypto/platform/telegram"
	"crypto/platform/utils"
	"os"
	"sync"
)

func main() {
	logger := utils.NewLogger()

	db, err := db.NewPostgresDBWithIDGen(DBURL())
	if err != nil {
		logger.Error("failed to connect to db", "error", err)
		return
	}
	defer db.Close()
	if err := db.Migrate(); err != nil {
		logger.Error("failed do migrations", "error", err)
	}

	a := app.NewApp(logger, db)

	run(a)
}

func DBURL() string {
	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return "postgresql://user:password@localhost:5432/dbname?sslmode=disable"
	}
	return url
}

func run(a *app.App) {
	var wg sync.WaitGroup

	for _, f := range toRun() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f(a)
		}()
	}

	wg.Wait()
}

func toRun() [3]func(*app.App) {
	updated := make(chan struct{}, 1)

	return [3]func(*app.App){
		func(a *app.App) { startCollecting(a, updated) },
		func(a *app.App) { startMonitoring(a, updated) },
		startTgBot,
	}
}

func startCollecting(a *app.App, updated chan<- struct{}) {
	c := collectors.NewRateCollector(a, updated)
	toRun := func() error { return c.Run() }
	if err := utils.LogProcess(*a.Logger, "collecting", toRun); err != nil {
		a.Logger.Error("failed collect logs", "error", err)
	}
}

func startMonitoring(a *app.App, updated <-chan struct{}) {
	m := monitoring.NewMonitoring(a, updated)
	toRun := func() error { return m.Run() }
	if err := utils.LogProcess(*a.Logger, "collecting", toRun); err != nil {
		a.Logger.Error("failed collect logs", "error", err)
	}
}

func startTgBot(a *app.App) {
	b := telegram.NewBotRunner(a)
	toRun := func() error { return b.Run() }
	if err := utils.LogProcess(*a.Logger, "tg bot", toRun); err != nil {
		a.Logger.Error("failed run Telegram Bot", "error", err)
	}
}
