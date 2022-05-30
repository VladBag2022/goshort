package shortener

import (
	"math/rand"
	"net/url"

	"github.com/google/uuid"
)

func Shorten(_ *url.URL) (string, error) {
	return RandomString(10), nil
}

func Register() string {
	return uuid.New().String()
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}