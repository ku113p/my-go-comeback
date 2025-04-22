package models

import "fmt"

type NotExistsError struct {
	ModelName ModelName
	ID        ObjectID
}

func NewNotExistsError(name ModelName, id ObjectID) *NotExistsError {
	return &NotExistsError{name, id}
}

func (e *NotExistsError) Error() string {
	return fmt.Sprintf("%s(id=%v) not exists", e.ModelName, e.ID)
}

type AlreadyExistsError struct {
	ModelName ModelName
	ID        ObjectID
}

func NewAlreadyExistsError(name ModelName, id ObjectID) *AlreadyExistsError {
	return &AlreadyExistsError{name, id}
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s(id=%v) already exists", e.ModelName, e.ID)
}
