package storage

import (
	"net/url"
)

type ShortURL struct {
	id int
	Origin url.URL
	Result url.URL
}