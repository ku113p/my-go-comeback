package db

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
)

const storageFp = "./storage.json"

type Task struct {
	Id          uuid.UUID `json:"id"`
	Time        time.Time `json:"time"`
	Done        bool      `json:"done"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
}

type TaskBuilder struct {
	id          uuid.UUID
	time        time.Time
	name        string
	description string
}

func NewTaskBuilder() *TaskBuilder {
	return &TaskBuilder{}
}

func (t *TaskBuilder) WithSomeId() *TaskBuilder {
	id, err := uuid.NewV7()
	if err != nil {
		panic("No way! It's happened.")
	}
	return t.withId(id)
}

func (t *TaskBuilder) withId(id uuid.UUID) *TaskBuilder {
	t.id = id
	return t
}

func (t *TaskBuilder) WithTime(time time.Time) *TaskBuilder {
	t.time = time
	return t
}

func (t *TaskBuilder) WithName(n string) *TaskBuilder {
	t.name = n
	return t
}

func (t *TaskBuilder) WithDescription(d string) *TaskBuilder {
	t.description = d
	return t
}

func (t *TaskBuilder) Build() *Task {
	return &Task{
		Id:          t.id,
		Time:        t.time,
		Name:        t.name,
		Description: t.description,
	}
}

func GetDataFromFs() (map[string]*Task, error) {
	jsonFile, err := os.Open(storageFp)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var d map[string]*Task

	if err := json.Unmarshal(byteValue, &d); err != nil {
		return nil, err
	}

	return d, nil
}

func saveDataToFs(d map[string]*Task) error {
	byteValue, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}

	jsonFile, err := os.OpenFile(storageFp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonFile.Write(byteValue)

	return nil
}
