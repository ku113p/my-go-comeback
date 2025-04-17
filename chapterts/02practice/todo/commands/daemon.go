package commands

import (
	"os"
	"todo/cli/db"
)

func daemon(src string) {
	wd, _ := os.Getwd()
	logger.Info("Daemon started", "wd", wd, "src", src)

	s := db.GetStorage()
	s.StartSaveEveryMinute()
	s.MonitorOperationsEvery10Seconds(src)

	select {}
}
