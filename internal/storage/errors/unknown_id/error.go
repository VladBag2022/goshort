package unknown_id

func New(id int) error {
	return &UnknownIDError{id}
}

type UnknownIDError struct {
	id int
}

func (e *UnknownIDError) Error() string {
	return string(rune(e.id))
}