package storage

import (
	"context"
	"net/url"
)

type Repository interface {
	Shorten(
		ctx context.Context,
		origin url.URL,
	) (*ShortURL, error)

	Restore(
		ctx context.Context,
		id int,
	) (*ShortURL, error)
}