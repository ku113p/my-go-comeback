package commands

import (
	"todo/cli/db"

	"github.com/google/uuid"
)

func listTasks() []*db.Task {
	return db.GetStorage().ListTasks()
}

func getTask(id string) (*db.Task, bool) {
	return db.GetStorage().GetTask(id)
}

func deleteTask(id string) error {
	s := db.GetStorage()
	if err := s.DeleteTask(id); err != nil {
		return err
	}
	if err := s.Save(); err != nil {
		return err
	}
	return nil
}

func markDoneTask(id string) error {
	s := db.GetStorage()
	if err := s.MarkDone(id); err != nil {
		return err
	}
	if err := s.Save(); err != nil {
		return err
	}
	return nil
}

func newTask(name, desc string) (*uuid.UUID, error) {
	t := db.NewTaskBuilder(db.UuidIdGenerator).WithName(name).WithDescription(desc).Build()

	s := db.GetStorage()
	if err := s.AddTask(t); err != nil {
		return nil, err
	}
	if err := s.Save(); err != nil {
		return nil, err
	}

	return &t.ID, nil
}
