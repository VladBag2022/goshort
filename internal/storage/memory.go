package storage

import (
	"context"
	"github.com/VladBag2022/goshort/internal/storage/errors/unknown_id"
	"net/url"
	"sync"
)

type MemoryRepository struct {
	urls 	  map[int]ShortURL
	nextId 	  int
	lock  	  *sync.RWMutex
	shortenFn func(url.URL) (url.URL, error)
}

func NewMemoryRepository(shortenFn func(url.URL) (url.URL, error)) *MemoryRepository {
	return &MemoryRepository{
		urls:     	map[int]ShortURL{},
		nextId:	  	0,
		lock:     	&sync.RWMutex{},
		shortenFn: 	shortenFn,
	}
}

func (m MemoryRepository) Shorten(_ context.Context, origin url.URL) (*ShortURL, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.shorten(origin)
}

func (m MemoryRepository) Restore(_ context.Context, id int) (*ShortURL, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.restore(id)
}

func (m MemoryRepository) shorten(origin url.URL) (*ShortURL, error) {
	result, err := m.shortenFn(origin)
	if err != nil {
		return nil, err
	}
	shortURL := ShortURL{
		id: m.nextId,
		Origin: origin,
		Result: result,
	}
	m.urls[shortURL.id] = shortURL
	m.nextId += 1
	return &shortURL, nil
}

func (m MemoryRepository) restore(id int) (*ShortURL, error) {
	shortURL, ok := m.urls[id]
	if !ok {
		return nil, unknown_id.New(id)
	}
	return &shortURL, nil
}