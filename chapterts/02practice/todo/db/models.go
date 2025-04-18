package db

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

const storageFp = "./storage.json"

var appFs afero.Fs

func init() {
	appFs = afero.NewOsFs()
}

type Task struct {
	ID          uuid.UUID `json:"id"`
	Time        time.Time `json:"time"`
	Done        bool      `json:"done"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
}

type TaskBuilder struct {
	genId       func() uuid.UUID
	id          uuid.UUID
	time        time.Time
	name        string
	description string
}

func NewTaskBuilder(genId func() uuid.UUID) *TaskBuilder {
	return &TaskBuilder{genId: genId}
}

func UuidIdGenerator() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		logger.Error("Failed gen v7 Uuid", "error", err)
		return uuid.New()
	}
	return id
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
	t.withId(t.genId())

	return &Task{
		ID:          t.id,
		Time:        t.time,
		Name:        t.name,
		Description: t.description,
	}
}

func getDataFromFs() (map[string]*Task, error) {
	jsonFile, err := appFs.Open(storageFp)
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

	jsonFile, err := appFs.OpenFile(storageFp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonFile.Write(byteValue)

	return nil
}
