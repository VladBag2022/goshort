package storage

import (
	"context"
	"net/url"
)

type Repository interface {
	Shorten(
		ctx 	context.Context,
		origin  *url.URL,
	) (string, error)

	Restore(
		ctx context.Context,
		id 	string,
	) (*url.URL, error)

	Load(ctx context.Context) error
	Dump(ctx context.Context) error

	Register(ctx context.Context) (string, error)
	Bind(
		ctx 	context.Context,
		urlID 	string,
		userID 	string,
	) error

	ShortenedList(
		ctx context.Context,
		id  string,
	) ([]string, error)

	Close() []error
}