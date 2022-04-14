package shortener

import "net/url"

func Shorten(u url.URL) (url.URL, error) {
	return u, nil
}