package storage

import "fmt"

type UnknownIDError struct {
	id string
}

func NewUnknownIDError(id string) error {
	return &UnknownIDError{id}
}

func (e *UnknownIDError) Error() string {
	return e.id
}

type NoCoolStorageError struct {
	repositoryType	string
}

func NewNoCoolStorageError(repositoryType	string) error {
	return &NoCoolStorageError{
		repositoryType:	repositoryType,
	}
}

func (e *NoCoolStorageError) Error() string {
	return fmt.Sprintf("CoolStorage was not provided during %s initialisation", e.repositoryType)
}