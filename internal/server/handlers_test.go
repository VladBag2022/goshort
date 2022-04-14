package server

import (
	"context"
	"github.com/VladBag2022/goshort/internal/shortener"
	"github.com/VladBag2022/goshort/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestServer_shorten(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response 	bool
	}
	tests := []struct {
		name    string
		content string
		want want
	}{
		{
			name: 		"positive test",
			content: 	"https://example.com",
			want: want{
				contentType: 	"text/plain; charset=utf-8",
				statusCode: 	201,
				response: 		true,
			},
		},
		{
			name: 		"negative test",
			content: 	"",
			want: want{
				contentType: 	"text/plain; charset=utf-8",
				statusCode: 	400,
				response: 		true,
			},
		},
	}

	r := storage.NewMemoryRepository(shortener.Shorten)
	s := New(r, "localhost", 8080)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.content))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(shortenHandler(&s))
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.response, len(userResult) > 0)
		})
	}
}

func TestServer_restore(t *testing.T) {
	type want struct {
		location    string
		statusCode  int
	}
	tests := []struct {
		name    string
		origin  string
		sameID  bool
		want want
	}{
		{
			name: 		"positive test",
			origin: 	"https://example.com",
			sameID:  	true,
			want: want{
				location: 	    "https://example.com",
				statusCode: 	307,
			},
		},
		{
			name: 		"negative test",
			origin: 	"https://example.com",
			sameID:  	false,
			want: want{
				location: 	    "https://example.com",
				statusCode: 	400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := storage.NewMemoryRepository(shortener.Shorten)
			u, err := url.Parse(tt.origin)
			require.NoError(t, err)
			id, err := r.Shorten(context.Background(), u)
			require.NoError(t, err)
			s := New(r, "localhost", 8080)

			requestURL := "/" + id
			if !tt.sameID {
				requestURL += "444"
			}

			request := httptest.NewRequest(http.MethodGet, requestURL, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(restoreHandler(&s))
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			if tt.sameID {
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			}

			_, err = ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}