package storage

func NewUnknownIDError(id string) error {
	return &UnknownIDError{id}
}

type UnknownIDError struct {
	id string
}

func (e *UnknownIDError) Error() string {
	return e.id
}