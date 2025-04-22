package models

type Table struct {
	name    ModelName
	storage map[ObjectID]ObjectsModel
}

func newTable(name ModelName) *Table {
	return &Table{name: name, storage: make(map[ObjectID]ObjectsModel)}
}

func (t Table) List() ([]ObjectsModel, error) {
	result := make([]ObjectsModel, 0)

	for _, m := range t.storage {
		result = append(result, m)
	}

	return result, nil
}

func (t Table) Get(id ObjectID) (ObjectsModel, error) {
	obj, ok := t.storage[id]
	if !ok {
		return nil, NewNotExistsError(t.name, id)
	}
	return obj, nil
}

func (t Table) Delete(id ObjectID) error {
	_, ok := t.storage[id]
	if !ok {
		return NewNotExistsError(t.name, id)
	}

	delete(t.storage, id)

	return nil
}

func (t Table) Update(obj ObjectsModel) error {
	_, ok := t.storage[obj.getID()]
	if !ok {
		return NewNotExistsError(t.name, obj.getID())
	}

	t.storage[obj.getID()] = obj

	return nil
}

func (t Table) Create(obj ObjectsModel) error {
	_, ok := t.storage[obj.getID()]
	if !ok {
		t.storage[obj.getID()] = obj
		return nil
	}

	return NewAlreadyExistsError(t.name, obj.getID())
}
