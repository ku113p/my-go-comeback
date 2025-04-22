package models

import (
	"encoding/json"
	"io"

	"github.com/google/uuid"
)

type ObjectID uuid.UUID

func (id ObjectID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id).String())
}

func (id *ObjectID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	*id = ObjectID(parsed)
	return nil
}

type IDGenerator interface {
	Generate() ObjectID
}

type ModelName string

type ModelsRepositry struct {
	db          map[ModelName]objectsStorage
	idGenerator IDGenerator
}

type objectsStorage interface {
	List() ([]ObjectsModel, error)
	Create(ObjectsModel) error
	Get(ObjectID) (ObjectsModel, error)
	Update(ObjectsModel) error
	Delete(ObjectID) error
}

type ObjectsModel interface {
	getID() ObjectID
	SetID(*ObjectID)
}

var ModelsToRegister = []ModelName{noteModelName, authorModelName}

func NewModelsRepository(idGenerator IDGenerator, models []ModelName) *ModelsRepositry {
	db := make(map[ModelName]objectsStorage, 0)

	for _, m := range models {
		db[m] = newTable(m)
	}

	return &ModelsRepositry{
		db:          db,
		idGenerator: idGenerator,
	}
}

type uuidIdGenerator struct{}

func NewUuidGenerator() IDGenerator {
	return uuidIdGenerator{}
}

func (u uuidIdGenerator) Generate() ObjectID {
	id, err := uuid.NewV7()
	if err != nil {
		panic("impossible")
	}
	return ObjectID(id)
}

var ModelParsers = map[ModelName]func(io.Reader) (ObjectsModel, error){
	noteModelName:   noteFromBytes,
	authorModelName: authorFromBytes,
}
