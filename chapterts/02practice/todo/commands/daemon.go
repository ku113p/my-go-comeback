package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"todo/cli/db"

	"github.com/google/uuid"
)

func daemon(src string) {
	wd, _ := os.Getwd()
	logger.Info("Daemon started", "wd", wd, "src", src)

	s := db.GetStorage()
	s.StartSaveEveryMinute()
	monitorOperations(src, s)

	select {}
}

func monitorOperations(src string, s *db.Storage) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			<-ticker.C
			c := dirOperations(src, s)
			logger.Info("made operations", "counter", c)
		}
	}()
}

func dirOperations(src string, s *db.Storage) uint {
	var counter uint = 0

	entries, err := os.ReadDir(src)
	if err != nil {
		msg := fmt.Sprintf("Failed read dir %v", src)
		logger.Error(msg, "error", err)
		return counter
	}

	filePaths := []string{}
	for _, item := range entries {
		if !item.IsDir() {
			filePaths = append(filePaths, filepath.Join(src, item.Name()))
		}
	}

	for _, fp := range filePaths {
		if err := makeOperation(fp, s); err != nil {
			logger.Error("Failed makeOperation", "fp", fp, "error", err)
		} else {
			counter++
			if err := os.Remove(fp); err != nil {
				logger.Error("Failed delete operation file", "fp", fp)
			}
		}
	}

	return counter
}

const (
	deleteOpName = "delete"
	markOpName   = "mark"
	newOpName    = "new"
)

var allowedOp = map[string]bool{
	deleteOpName: true,
	markOpName:   true,
	newOpName:    true,
}

func makeOperation(src string, s *db.Storage) error {
	var o Operation

	operation := strings.SplitN(filepath.Base(src), "_", 2)[0]
	if !allowedOp[operation] {
		return fmt.Errorf("unknonw operation %v", operation)
	}

	jsonFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	switch operation {
	case deleteOpName:
		var d deleteOperation
		if err := json.Unmarshal(byteValue, &d); err != nil {
			return err
		}
		o = &d
	case markOpName:
		var m markDoneOperation
		if err := json.Unmarshal(byteValue, &m); err != nil {
			return err
		}
		o = &m
	case newOpName:
		var n newOperation
		if err := json.Unmarshal(byteValue, &n); err != nil {
			return err
		}
		o = &n
	}

	return o.make(s)
}

type Operation interface {
	make(s *db.Storage) error
}

type idOperation struct {
	Id uuid.UUID `json:"id"`
}

type deleteOperation struct {
	idOperation
}

func (d *deleteOperation) make(s *db.Storage) error {
	return s.DeleteTask(d.Id.String())
}

type markDoneOperation struct {
	idOperation
}

func (d *markDoneOperation) make(s *db.Storage) error {
	return s.MarkDone(d.Id.String())
}

type newOperation struct {
	Time        time.Time `json:"time"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
}

func (c *newOperation) make(s *db.Storage) error {
	return s.AddTask(
		db.NewTaskBuilder(db.UuidIdGenerator).
			WithName(c.Name).
			WithDescription(c.Description).
			WithTime(c.Time).
			Build(),
	)
}
