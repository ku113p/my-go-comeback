package db

import (
	"fmt"
	"os"
	"time"

	"log/slog"
)

var logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

type Storage struct {
	data   map[string]*Task
	locker chan any
}

func GetStorage() *Storage {
	data, err := GetDataFromFs()
	if err != nil {
		logger.Info("Storage not loaded from file system. New one will be created.", "error", err)
	}

	return newStorage(data)
}

func newStorage(data map[string]*Task) *Storage {
	if data == nil {
		data = map[string]*Task{}
	}

	return &Storage{
		data:   data,
		locker: make(chan any, 1),
	}
}

func (s *Storage) StartSaveEveryMinute() {
	s.startPeriodicalSave(time.Minute)
}

func (s *Storage) startPeriodicalSave(period time.Duration) {
	ticker := time.NewTicker(period)

	go func() {
		for {
			<-ticker.C
			if err := s.Save(); err != nil {
				logger.Error("Failed save tasks to Fs", "error", err)
			} else {
				logger.Info("Tasks state saved")
			}
		}
	}()
}

func (s *Storage) Save() error {
	return saveDataToFs(s.data)
}

func (s *Storage) ListTasks() []*Task {
	to_defer := s.borrowSpace()
	defer to_defer()

	v := make([]*Task, 0, len(s.data))
	for _, val := range s.data {
		v = append(v, val)
	}

	return v
}

func (s *Storage) GetTask(id string) (*Task, bool) {
	to_defer := s.borrowSpace()
	defer to_defer()

	t, ok := s.data[id]
	return t, ok
}

func (s *Storage) borrowSpace() func() {
	s.locker <- true
	return func() { <-s.locker }
}

func (s *Storage) AddTask(t *Task) error {
	to_defer := s.borrowSpace()
	defer to_defer()

	tid := t.ID.String()
	if _, exists := s.data[tid]; exists {
		return fmt.Errorf("task with id '%v' already exists", tid)
	}

	s.data[tid] = t

	return nil
}

func (s *Storage) DeleteTask(id string) error {
	to_defer := s.borrowSpace()
	defer to_defer()

	if _, exists := s.data[id]; exists {
		delete(s.data, id)
		return nil
	}

	return fmt.Errorf("task with id '%v' not exists", id)
}

func (s *Storage) MarkDone(id string) error {
	to_defer := s.borrowSpace()
	defer to_defer()

	if t, exists := s.data[id]; exists {
		if t.Done {
			return fmt.Errorf("task with id '%v' already marked done", id)
		}
		t.Done = true
		return nil
	}

	return fmt.Errorf("task with id '%v' not exists", id)
}
