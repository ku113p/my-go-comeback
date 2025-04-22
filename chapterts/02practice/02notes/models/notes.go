package models

import (
	"encoding/json"
	"io"
	"time"
)

const noteModelName ModelName = "note"

type Note struct {
	ID          *ObjectID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Created     *time.Time `json:"created"`
	AuthorId    *ObjectID  `json:"author_id"`
}

func noteFromBytes(r io.Reader) (ObjectsModel, error) {
	var note Note
	if err := json.NewDecoder(r).Decode(&note); err != nil {
		return nil, err
	}
	return &note, nil
}

func (n *Note) getID() ObjectID {
	return *n.ID
}

func (n *Note) SetID(id *ObjectID) {
	n.ID = id
}

func (n *Note) setDefaults() {
	if n.Created == nil {
		now := time.Now()
		n.Created = &now
	}
}
