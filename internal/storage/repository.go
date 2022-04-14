package storage

import (
	"context"
	"net/url"
)

type Repository interface {
	Shorten(
		ctx context.Context,
		origin *url.URL,
	) (string, error)

	Restore(
		ctx context.Context,
		id string,
	) (*url.URL, error)
}