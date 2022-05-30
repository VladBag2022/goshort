package misc

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

func Sign(key, value string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	dst := h.Sum(nil)

	return fmt.Sprintf("%s|%x", value, dst)
}

func Verify(key, msg string) (bool, string, error) {
	values := strings.Split(msg, "|")
	if len(values) != 2 {
		return false, "", errors.New("separator not found")
	}

	return Sign(key, values[0]) == msg, values[0], nil
}

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