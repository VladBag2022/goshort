package storage

import (
	"context"
	"net/url"
	"sync"
)

type MemoryRepository struct {
	urls      sync.Map
	shortenFn func(*url.URL) (string, error)
}

func NewMemoryRepository(shortenFn func(*url.URL) (string, error)) *MemoryRepository {
	return &MemoryRepository{
		shortenFn: shortenFn,
	}
}

func (m *MemoryRepository) Shorten(_ context.Context, origin *url.URL) (string, error) {
	var id = ""
	for id == "" {
		newID, err := m.shortenFn(origin)
			if err != nil{
			return "", err
		}
		_, ok := m.urls.Load(newID)
		if !ok {
			id = newID
		}
	}
	m.urls.Store(id, origin)
	return id, nil
}

func (m *MemoryRepository) Restore(_ context.Context, id string) (*url.URL, error) {
	origin, ok := m.urls.Load(id)
	if !ok {
		return nil, NewUnknownIDError(id)
	}
	return origin.(*url.URL), nil
}