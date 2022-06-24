package storage

import (
	"context"
	"fmt"
	"net/url"
	"sync"
)

type MemoryRepository struct {
	urls      		sync.Map
	userUrls 		sync.Map
	shortenFn 		func(*url.URL) (string, error)
	registerFn 		func() string
	coolStorage		*CoolStorage
}

func NewMemoryRepository(
	shortenFn 		func(*url.URL) (string, error),
	registerFn 		func() string,
) *MemoryRepository {
	return &MemoryRepository{
		shortenFn: 	shortenFn,
		registerFn: registerFn,
	}
}

func NewMemoryRepositoryWithCoolStorage(
	shortenFn 		func(*url.URL) (string, error),
	registerFn 		func() string,
	coolStorage 	*CoolStorage,
) *MemoryRepository {
	return &MemoryRepository{
		shortenFn: 		shortenFn,
		registerFn:     registerFn,
		coolStorage: 	coolStorage,
	}
}

func (m *MemoryRepository) Shorten(_ context.Context, origin *url.URL) (string, bool, error) {
	var id = ""
	for id == "" {
		newID, err := m.shortenFn(origin)
		if err != nil {
			return "", false, err
		}
		_, ok := m.urls.Load(newID)
		if !ok {
			id = newID
		}
	}
	m.urls.Store(id, origin)
	return id, true, nil
}

func (m *MemoryRepository) Restore(_ context.Context, id string) (*url.URL, bool, error) {
	origin, ok := m.urls.Load(id)
	if !ok {
		return nil, false, NewUnknownIDError(fmt.Sprintf("url: %s", id))
	}
	return origin.(*url.URL), false, nil
}

func (m *MemoryRepository) Load(_ context.Context) error {
	if m.coolStorage == nil {
		return NewNoCoolStorageError("MemoryRepository")
	}
	records, err := m.coolStorage.FetchRecords()
	if err != nil {
		return err
	}
	for _, record := range records {
		origin, err := url.Parse(record.Origin)
		if err != nil {
			continue
		}
		m.urls.Store(record.ID, origin)

		urls, ok := m.userUrls.Load(record.User)
		if !ok {
			urls = []string{}
		}
		urls = append(urls.([]string), record.ID)
		m.userUrls.Store(record.User, urls)
	}
	return nil
}

func (m *MemoryRepository) Dump(ctx context.Context) error {
	if m.coolStorage == nil {
		return NewNoCoolStorageError("MemoryRepository")
	}
	var records []*CoolStorageRecord
	m.userUrls.Range(func(userID, urlIDs interface{}) bool {
		for _, urlID := range urlIDs.([]string) {
			originURL, _, err := m.Restore(ctx, urlID)
			if err != nil {
				return true
			}
			records = append(records, &CoolStorageRecord{
				Origin: originURL.String(),
				ID: 	urlID,
				User: 	userID.(string),
			})
		}
		return true
	})
	return m.coolStorage.PutRecords(records)
}

func (m *MemoryRepository) Register(_ context.Context) (string, error) {
	var id = ""
	for id == "" {
		newID := m.registerFn()
		_, ok := m.urls.Load(newID)
		if !ok {
			id = newID
		}
	}
	m.userUrls.Store(id, []string{})
	return id, nil
}

func (m *MemoryRepository) Bind(
	_ 	 	context.Context,
	urlID 	string,
	userID 	string,
) error {
	urls, ok := m.userUrls.Load(userID)
	if !ok {
		return NewUnknownIDError(fmt.Sprintf("user: %s", userID))
	}
	_, ok = m.urls.Load(urlID)
	if !ok {
		return NewUnknownIDError(fmt.Sprintf("url: %s", userID))
	}

	for _, storedID := range urls.([]string) {
		if storedID == urlID {
			return nil
		}
	}

	m.userUrls.Store(userID, append(urls.([]string), urlID))
	return nil
}

func (m *MemoryRepository) ShortenedList(
	_ context.Context,
	id  string,
) ([]string, error) {
	urls, ok := m.userUrls.Load(id)
	if !ok {
		return []string{}, NewUnknownIDError(fmt.Sprintf("user: %s", id))
	}
	return urls.([]string), nil
}

func (m *MemoryRepository) Close() []error {
	var errs []error

	if m.coolStorage != nil {
		err := m.coolStorage.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (m *MemoryRepository) ShortenBatch(
	ctx			context.Context,
	origins		[]*url.URL,
	userID 		string,
) ([]string, error) {
	var ids []string
	for _, origin := range origins {
		id, _, err := m.Shorten(ctx, origin)
		if err != nil {
			return ids, err
		}
		err = m.Bind(ctx, id, userID)
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}