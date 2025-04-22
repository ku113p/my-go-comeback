package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"
)

const NoteModelName ModelName = "notes"

type Note struct {
	ID          *ObjectID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Created     time.Time  `json:"created"`
	AuthorId    *uuid.UUID `json:"author_id"`
}

func noteFromBytes(r io.Reader) (Model, error) {
	var note *Note
	if err := json.NewDecoder(r).Decode(&note); err != nil {
		return nil, err
	}
	return note, nil
}

func (n *Note) getID() ObjectID {
	return *n.ID
}

func (n *Note) SetID(id *ObjectID) {
	n.ID = id
}
