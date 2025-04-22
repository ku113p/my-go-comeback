package models

import (
	"encoding/json"
	"io"
)

const AuthorModelName ModelName = "authors"

type Author struct {
	ID         *ObjectID `json:"id"`
	Username   string    `json:"username"`
	Firstname  *string   `json:"firstname"`
	Secondname *string   `json:"secondname"`
}

func authorFromBytes(r io.Reader) (Model, error) {
	var author *Author
	if err := json.NewDecoder(r).Decode(author); err != nil {
		return nil, err
	}
	return author, nil
}

func (a *Author) getID() ObjectID {
	return *a.ID
}

func (a *Author) SetID(id *ObjectID) {
	a.ID = id
}
