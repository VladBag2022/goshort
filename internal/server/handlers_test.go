package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/VladBag2022/goshort/internal/misc"
	"github.com/VladBag2022/goshort/internal/storage"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = resp.Body.Close()
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestServer_shorten(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    bool
	}
	tests := []struct {
		name    string
		content string
		want    want
	}{
		{
			name:    "positive test",
			content: "https://example.com",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
				response:    true,
			},
		},
		{
			name:    "negative test",
			content: "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    true,
			},
		},
	}

	mem := storage.NewMemoryRepository(misc.Shorten, misc.UUID)
	defer mem.Close()

	c, err := NewConfig()
	if err != nil {
		require.NoError(t, err)
	}
	s := NewServer(mem, nil, c)

	r := router(s)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, content := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(tt.content))
			err := response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, len(content) > 0)
		})
	}
}

func TestServer_api_shorten(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    bool
	}
	tests := []struct {
		name    string
		content string
		want    want
	}{
		{
			name:    "positive test",
			content: "{\"url\":\"https://example.com\"}",
			want: want{
				contentType: "application/json",
				statusCode:  201,
				response:    true,
			},
		},
		{
			name:    "negative test",
			content: "{\"link\":\"https://example.com\"}",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				response:    true,
			},
		},
	}

	mem := storage.NewMemoryRepository(misc.Shorten, misc.UUID)
	defer mem.Close()
	c, err := NewConfig()
	if err != nil {
		require.NoError(t, err)
	}
	s := NewServer(mem, nil, c)

	r := router(s)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, content := testRequest(t, ts, http.MethodPost, "/api/shorten", strings.NewReader(tt.content))
			err := response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, len(content) > 0)
		})
	}
}

func TestServer_restore(t *testing.T) {
	type want struct {
		location   string
		statusCode int
	}
	tests := []struct {
		name   string
		origin string
		sameID bool
		want   want
	}{
		{
			name:   "positive test",
			origin: "https://example.com",
			sameID: true,
			want: want{
				location:   "https://example.com",
				statusCode: 307,
			},
		},
		{
			name:   "negative test",
			origin: "https://example.com",
			sameID: false,
			want: want{
				location:   "https://example.com",
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := storage.NewMemoryRepository(misc.Shorten, misc.UUID)
			defer mem.Close()

			u, err := url.Parse(tt.origin)
			require.NoError(t, err)
			id, _, err := mem.Shorten(context.Background(), u)
			require.NoError(t, err)

			c, err := NewConfig()
			if err != nil {
				require.NoError(t, err)
			}
			s := NewServer(mem, nil, c)
			r := router(s)
			ts := httptest.NewServer(r)
			defer ts.Close()

			requestURL := "/" + id
			if !tt.sameID {
				requestURL += "GARBAGE"
			}

			t.Logf("ID: %s", id)
			t.Logf("Request: %s", requestURL)

			response, content := testRequest(t, ts, http.MethodGet, requestURL, nil)
			err = response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			if tt.sameID {
				assert.Equal(t, tt.want.location, response.Header.Get("Location"))
			}

			t.Logf("Body: %s", content)
		})
	}
}
