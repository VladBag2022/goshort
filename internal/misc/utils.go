package misc

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"sync"

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

func UUID() string {
	return uuid.New().String()
}

func RandomString(n int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func FanIn(inputChs ...chan interface{}) chan interface{} {
	outCh := make(chan interface{})

	go func() {
		wg := &sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)

			go func(inputCh chan interface{}) {
				defer wg.Done()
				for item := range inputCh {
					outCh <- item
				}
			}(inputCh)
		}

		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func FanOut(inputCh chan interface{}, n int) []chan interface{} {
	chs := make([]chan interface{}, 0, n)
	for i := 0; i < n; i++ {
		ch := make(chan interface{})
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan interface{}) {
			for _, ch := range chs {
				close(ch)
			}
		}(chs)

		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}

			num, ok := <-inputCh
			if !ok {
				return
			}

			ch := chs[i]
			ch <- num
		}
	}()

	return chs
}
