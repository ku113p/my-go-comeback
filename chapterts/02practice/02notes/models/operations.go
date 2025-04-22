package models

type RepositoryOperation struct {
	Name       ModelName
	repository *ModelsRepositry
}

func NewRepositoryOperation(name ModelName, repository *ModelsRepositry) *RepositoryOperation {
	return &RepositoryOperation{name, repository}
}

func (op *RepositoryOperation) table() Storage {
	return op.repository.db[op.Name]
}

func (op *RepositoryOperation) List() ([]Model, error) {
	return op.table().List()
}

func (op *RepositoryOperation) Create(m Model) error {
	id := op.repository.idGenerator.Generate()
	m.SetID(&id)
	return op.table().Create(m)
}

func (op *RepositoryOperation) Get(id ObjectID) (Model, error) {
	return op.table().Get(id)
}

func (op *RepositoryOperation) Update(m Model) error {
	return op.table().Update(m)
}

func (op *RepositoryOperation) Delete(id ObjectID) error {
	return op.table().Delete(id)
}
