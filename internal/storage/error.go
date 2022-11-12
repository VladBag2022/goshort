package storage

import "fmt"

type UnknownIDError struct {
	ID string
}

func NewUnknownIDError(id string) error {
	return &UnknownIDError{id}
}

func (e *UnknownIDError) Error() string {
	return e.ID
}

type noCoolStorageError struct {
	repositoryType string
}

func NewNoCoolStorageError(repositoryType string) error {
	return &noCoolStorageError{
		repositoryType: repositoryType,
	}
}

func (e *noCoolStorageError) Error() string {
	return fmt.Sprintf("CoolStorage was not provided during %s initialisation", e.repositoryType)
}
