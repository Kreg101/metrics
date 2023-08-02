package handler

import (
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/Kreg101/metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux(t *testing.T) {
	tt := []struct {
		name     string
		param    Repository
		expected *Mux
	}{
		{
			name:     "nil storage",
			param:    nil,
			expected: &Mux{storage: nil, log: logger.Default()},
		},
		{
			name:     "default storage",
			param:    &storage.Storage{},
			expected: &Mux{storage: &storage.Storage{}, log: logger.Default()},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mux := NewMux(tc.param, nil)
			assert.Equal(t, tc.expected, mux)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestMux_Router(t *testing.T) {
	s, err := storage.NewStorage("", 0, false, false, nil)
	require.NoError(t, err)

	mux := NewMux(s, nil)
	ts := httptest.NewServer(mux.Router())
	defer ts.Close()
	type response struct {
		statusCode  int
		contentType string
		body        string
	}
	tt := []struct {
		name   string
		url    string
		method string
		want   response
	}{
		{
			name:   "main page",
			url:    "/",
			method: http.MethodGet,
			want:   response{http.StatusOK, "text/html", ""},
		},
		{
			name:   "correct update counter request",
			url:    "/update/counter/x/10",
			method: http.MethodPost,
			want:   response{http.StatusOK, "", ""},
		},
		{
			name:   "correct update gauge request",
			url:    "/update/gauge/y/1.23",
			method: http.MethodPost,
			want:   response{http.StatusOK, "", ""},
		},
		{
			name:   "no metric name in update request #1",
			url:    "/update/counter//10",
			method: http.MethodPost,
			want:   response{http.StatusNotFound, "", ""},
		},
		{
			name:   "no metric name in update request #2",
			url:    "/update/counter",
			method: http.MethodPost,
			want:   response{http.StatusNotFound, "text/plain; charset=utf-8", "404 page not found\n"},
		},
		{
			name:   "invalid counter type value",
			url:    "/update/counter/x/abc",
			method: http.MethodPost,
			want:   response{http.StatusBadRequest, "", ""},
		},
		{
			name:   "invalid gauge type value",
			url:    "/update/counter/x/abc",
			method: http.MethodPost,
			want:   response{http.StatusBadRequest, "", ""},
		},
		{
			name:   "invalid metric type",
			url:    "/update/counte/x/abc",
			method: http.MethodPost,
			want:   response{http.StatusBadRequest, "", ""},
		},
		{
			name:   "single metric request",
			url:    "/value/counter/x",
			method: http.MethodGet,
			want:   response{http.StatusOK, "text/plain", "10"},
		},
		{
			name:   "invalid metric type request",
			url:    "/value/cor/x",
			method: http.MethodGet,
			want:   response{http.StatusNotFound, "", ""},
		},
		{
			name:   "no metric in storage",
			url:    "/value/counter/z",
			method: http.MethodGet,
			want:   response{http.StatusNotFound, "", ""},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tc.method, tc.url)
			defer resp.Body.Close()
			assert.Equal(t, tc.want.statusCode, resp.StatusCode)
			assert.Equal(t, tc.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tc.want.body, get)
		})
	}
}
