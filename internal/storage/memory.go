package storage

import (
	"context"
	"net/url"
	"sync"
)

type MemoryRepository struct {
	urls      		sync.Map
	shortenFn 		func(*url.URL) (string, error)
	coolStorage		*CoolStorage
}

func NewMemoryRepository(shortenFn func(*url.URL) (string, error)) *MemoryRepository {
	return &MemoryRepository{
		shortenFn: 		shortenFn,
	}
}

func NewMemoryRepositoryWithCoolStorage(shortenFn func(*url.URL) (string, error),
	coolStorage *CoolStorage) *MemoryRepository {
	return &MemoryRepository{
		shortenFn: 		shortenFn,
		coolStorage: 	coolStorage,
	}
}

func (m *MemoryRepository) Shorten(_ context.Context, origin *url.URL) (string, error) {
	var id = ""
	for id == "" {
		newID, err := m.shortenFn(origin)
			if err != nil {
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

func (m *MemoryRepository) Load(_ context.Context) error {
	if m.coolStorage == nil {
		return NewNoCoolStorageError("MemoryRepository")
	}
	records, err := m.coolStorage.Load()
	if err != nil {
		return err
	}
	for _, record := range records {
		origin, err := url.Parse(record.Origin)
		if err != nil {
			continue
		}
		m.urls.Store(record.ID, origin)
	}
	return nil
}

func (m *MemoryRepository) Dump(_ context.Context) error {
	if m.coolStorage == nil {
		return NewNoCoolStorageError("MemoryRepository")
	}
	var records []*CoolStorageRecord
	m.urls.Range(func(id, origin interface{}) bool {
		originURL := origin.(*url.URL)
		records = append(records, &CoolStorageRecord{
			Origin: originURL.String(),
			ID: 	id.(string),
		})
		return true
	})
	return m.coolStorage.Dump(records)
}