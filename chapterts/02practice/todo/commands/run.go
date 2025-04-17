package commands

import (
	"flag"
	"log/slog"
	"os"
	"strings"
)

var logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

func Run() {
	var fDaemonSrc string
	flag.StringVar(&fDaemonSrc, "daemon", "", "operations source dir for daemon")

	var fList bool
	flag.BoolVar(&fList, "list", false, "list tasks")

	var fGet string
	flag.StringVar(&fGet, "get", "", "get task by id")

	var fDel string
	flag.StringVar(&fDel, "del", "", "delete task by id")

	var fDone string
	flag.StringVar(&fDone, "done", "", "mark done task by id")

	var fNew string
	flag.StringVar(&fNew, "new", "", "create new task by \"<name>|<description>\"")

	flag.Parse()

	switch {
	case fDaemonSrc != "":
		daemon(fDaemonSrc)
	case fList:
		logger.Info("task", "tasks", listTasks())
	case fGet != "":
		if t, ok := getTask(fGet); ok {
			logger.Info("got task", "task", t)
		} else {
			logger.Warn("task not exists", "id", fGet)
		}
	case fDel != "":
		if err := deleteTask(fDel); err == nil {
			logger.Info("task deleted", "id", fDel)
		} else {
			logger.Warn("task delete error", "error", err)
		}
	case fDone != "":
		if err := markDoneTask(fDone); err == nil {
			logger.Info("task marked done", "id", fDone)
		} else {
			logger.Warn("task mark done error", "error", err)
		}
	case fNew != "":
		rows := strings.SplitN(fNew, "|", 2)
		if len(rows) == 2 {
			if id, err := newTask(rows[0], rows[1]); err == nil {
				logger.Info("task created", "id", id)
			} else {
				logger.Warn("task create error", "error", err)
			}
		} else {
			logger.Warn("invalid new task format", "input", fNew, "expected", "<name>|<description>")
		}
	default:
		logger.Warn("Empty call")
	}
}
