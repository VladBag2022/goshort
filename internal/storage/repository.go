package storage

import (
	"context"
	"net/url"
)

type Repository interface {
	Shorten(
		ctx context.Context,
		origin *url.URL,
	) (id string, created bool, err error)

	Restore(
		ctx context.Context,
		id string,
	) (origin *url.URL, deleted bool, err error)

	Delete(
		ctx context.Context,
		userID string,
		urlIDs []string,
	) error

	Load(ctx context.Context) error
	Dump(ctx context.Context) error

	Register(ctx context.Context) (id string, err error)
	Bind(
		ctx context.Context,
		urlID string,
		userID string,
	) error

	ShortenedList(
		ctx context.Context,
		id string,
	) (ids []string, err error)

	ShortenBatch(
		ctx context.Context,
		origins []*url.URL,
		userID string,
	) (ids []string, err error)

	Close() []error
}
