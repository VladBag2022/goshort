package http

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
	"github.com/VladBag2022/goshort/internal/server"
	"github.com/VladBag2022/goshort/internal/storage"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path, ip string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	if len(ip) > 0 {
		req.Header.Set("X-Real-IP", ip)
	}

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

func TestServer_stats(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name          string
		trustedSubnet string
		userIP        string
		userUrls      map[string][]string
		want          want
	}{
		{
			name:          "positive test, 2 users, 6 different urls, 6 total urls",
			trustedSubnet: "10.0.0.0/8",
			userIP:        "10.10.10.10",
			userUrls: map[string][]string{
				"john": {"http://url1", "http://url2", "http://url3"},
				"mark": {"http://url4", "http://url5", "http://url6"},
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				response:    "{\"urls\":6,\"users\":2}",
			},
		},
		{
			name:          "positive test, 2 users, 4 different urls, 6 total urls",
			trustedSubnet: "10.0.0.0/8",
			userIP:        "10.10.10.10",
			userUrls: map[string][]string{
				"john": {"http://url1", "http://url2", "http://url3"},
				"mark": {"http://url3", "http://url2", "http://url6"},
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				response:    "{\"urls\":6,\"users\":2}",
			},
		},
		{
			name:          "positive test, 2 users, 3 different urls, 6 total urls",
			trustedSubnet: "10.0.0.0/8",
			userIP:        "10.10.10.10",
			userUrls: map[string][]string{
				"john": {"http://url1", "http://url2", "http://url3"},
				"mark": {"http://url1", "http://url2", "http://url3"},
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				response:    "{\"urls\":6,\"users\":2}",
			},
		},
		{
			name:          "positive test, no users, no urls",
			trustedSubnet: "10.0.0.0/8",
			userIP:        "10.10.10.10",
			userUrls:      map[string][]string{},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				response:    "{\"urls\":0,\"users\":0}",
			},
		},
		{
			name:          "positive test, 2 users, no urls",
			trustedSubnet: "10.0.0.0/8",
			userIP:        "10.10.10.10",
			userUrls: map[string][]string{
				"john": {},
				"mark": {},
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				response:    "{\"urls\":0,\"users\":2}",
			},
		},
		{
			name:          "negative test, wrong address",
			trustedSubnet: "10.0.0.0/24",
			userIP:        "10.10.10.10",
			userUrls:      map[string][]string{},
			want: want{
				contentType: "",
				statusCode:  http.StatusForbidden,
				response:    "",
			},
		},
		{
			name:          "negative test, no trusted subnet",
			trustedSubnet: "",
			userIP:        "10.10.10.10",
			userUrls:      map[string][]string{},
			want: want{
				contentType: "",
				statusCode:  http.StatusForbidden,
				response:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := storage.NewMemoryRepository(misc.Shorten, misc.UUID)
			defer mem.Close()

			for _, urls := range tt.userUrls {
				_, err := mem.Register(context.Background())
				require.NoError(t, err)

				for _, tURL := range urls {
					u, uErr := url.Parse(tURL)
					require.NoError(t, uErr)
					_, _, uErr = mem.Shorten(context.Background(), u)
					require.NoError(t, uErr)
				}
			}

			c := server.NewConfig()
			c.TrustedSubnet = tt.trustedSubnet
			a := server.NewServer(mem, nil, c)
			s := NewServer(&a)

			r := router(s)
			ts := httptest.NewServer(r)
			defer ts.Close()

			response, content := testRequest(t, ts, http.MethodGet, "/api/internal/stats", tt.userIP, nil)
			err := response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, content)
		})
	}
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

	c := server.NewConfig()
	a := server.NewServer(mem, nil, c)
	s := NewServer(&a)

	r := router(s)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, content := testRequest(t, ts, http.MethodPost, "/", "", strings.NewReader(tt.content))
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
	c := server.NewConfig()
	a := server.NewServer(mem, nil, c)
	s := NewServer(&a)

	r := router(s)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, content := testRequest(t, ts, http.MethodPost, "/api/shorten", "", strings.NewReader(tt.content))
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

			c := server.NewConfig()
			a := server.NewServer(mem, nil, c)
			s := NewServer(&a)
			r := router(s)
			ts := httptest.NewServer(r)
			defer ts.Close()

			requestURL := "/" + id
			if !tt.sameID {
				requestURL += "GARBAGE"
			}

			t.Logf("ID: %s", id)
			t.Logf("Request: %s", requestURL)

			response, content := testRequest(t, ts, http.MethodGet, requestURL, "", nil)
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

func TestServer_delete(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name   string
		data   string
		want   want
	}{
		{
			name:   "positive test",
			data: "[\"123\",\"456\"]",
			want: want{
				statusCode: http.StatusAccepted,
			},
		},
		{
			name:   "negative test - wrong input format",
			data: "\"123\",\"456\"",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	mem := storage.NewMemoryRepository(misc.Shorten, misc.UUID)
	defer mem.Close()

	c := server.NewConfig()
	a := server.NewServer(mem, nil, c)
	s := NewServer(&a)
	r := router(s)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := testRequest(t, ts, http.MethodDelete, "/api/user/urls", "", strings.NewReader(tt.data))
			err := response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}

func TestServer_shortened_list(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name   string
		data   string
		want   want
	}{
		{
			name:   "positive test",
			data: "[\"123\",\"456\"]",
			want: want{
				statusCode: http.StatusAccepted,
			},
		},
		{
			name:   "negative test - wrong input format",
			data: "\"123\",\"456\"",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	mem := storage.NewMemoryRepository(misc.Shorten, misc.UUID)
	defer mem.Close()

	c := server.NewConfig()
	a := server.NewServer(mem, nil, c)
	s := NewServer(&a)
	r := router(s)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, _ := testRequest(t, ts, http.MethodDelete, "/api/user/urls", "", strings.NewReader(tt.data))
			err := response.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
