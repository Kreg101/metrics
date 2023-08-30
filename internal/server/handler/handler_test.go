package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/inmemstore"
	"github.com/Kreg101/metrics/internal/server/logger"
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
			name:     "nil inmemstore",
			param:    nil,
			expected: &Mux{storage: nil, log: logger.Default()},
		},
		{
			name:     "default inmemstore",
			param:    &inmemstore.InMemStorage{},
			expected: &Mux{storage: &inmemstore.InMemStorage{}, log: logger.Default()},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mux := NewMux(tc.param, nil, "")
			assert.Equal(t, tc.expected, mux)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body *metric.Metric) (*http.Response, string) {
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, ts.URL+path, nil)
		require.NoError(t, err)
	} else {

		js, err := json.Marshal(*body)
		require.NoError(t, err)

		var b bytes.Buffer
		w := gzip.NewWriter(&b)

		_, err = w.Write(js)
		require.NoError(t, err)

		err = w.Close()
		require.NoError(t, err)

		req, err = http.NewRequest(method, ts.URL+path, &b)
		req.Header.Set("Accept-Encoding", "gzip")
		req.Header.Set("Content-Encoding", "gzip")
		require.NoError(t, err)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var respBody string
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		defer gz.Close()

		var b bytes.Buffer
		_, err = b.ReadFrom(gz)
		require.NoError(t, err)

		respBody = b.String()
	} else {
		x, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		respBody = string(x)
	}
	return resp, respBody
}

func TestMux_Router(t *testing.T) {
	counter := int64(10)
	gauge := 1.2345
	s, err := inmemstore.NewInMemStorage("", 0, false, nil)
	require.NoError(t, err)

	mux := NewMux(s, nil, "")
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
		body   *metric.Metric
		want   response
	}{
		{
			name:   "main page",
			url:    "/",
			method: http.MethodGet,
			body:   nil,
			want:   response{http.StatusOK, "text/html", ""},
		},
		{
			name:   "correct update counter request",
			url:    "/update/counter/x/10",
			method: http.MethodPost,
			body:   nil,
			want:   response{http.StatusOK, "", ""},
		},
		{
			name:   "correct update gauge request",
			url:    "/update/gauge/y/1.23",
			method: http.MethodPost,
			body:   nil,
			want:   response{http.StatusOK, "", ""},
		},
		{
			name:   "no metric name in update request #1",
			url:    "/update/counter//10",
			method: http.MethodPost,
			body:   nil,
			want:   response{http.StatusNotFound, "", ""},
		},
		{
			name:   "no metric name in update request #2",
			url:    "/update/counter",
			method: http.MethodPost,
			body:   nil,
			want:   response{http.StatusNotFound, "text/plain; charset=utf-8", "404 page not found\n"},
		},
		{
			name:   "invalid counter type value",
			url:    "/update/counter/x/abc",
			method: http.MethodPost,
			body:   nil,
			want:   response{http.StatusBadRequest, "", ""},
		},
		{
			name:   "invalid gauge type value",
			url:    "/update/counter/x/abc",
			method: http.MethodPost,
			body:   nil,
			want:   response{http.StatusBadRequest, "", ""},
		},
		{
			name:   "invalid metric type",
			url:    "/update/counte/x/abc",
			method: http.MethodPost,
			body:   nil,
			want:   response{http.StatusBadRequest, "", ""},
		},
		{
			name:   "single metric request",
			url:    "/value/counter/x",
			method: http.MethodGet,
			body:   nil,
			want:   response{http.StatusOK, "text/plain", "10"},
		},
		{
			name:   "invalid metric type request",
			url:    "/value/cor/x",
			method: http.MethodGet,
			body:   nil,
			want:   response{http.StatusNotFound, "", ""},
		},
		{
			name:   "no metric in inmemstore",
			url:    "/value/counter/z",
			method: http.MethodGet,
			body:   nil,
			want:   response{http.StatusNotFound, "", ""},
		},
		{
			name:   "update counter metric with body",
			url:    "/update/",
			method: http.MethodPost,
			body:   &metric.Metric{ID: "key", MType: "counter", Delta: &counter},
			want: response{http.StatusOK, "application/json",
				"{\"id\":\"key\",\"type\":\"counter\",\"delta\":10}\n"},
		},
		{
			name:   "update gauge metric with body",
			url:    "/update/",
			method: http.MethodPost,
			body:   &metric.Metric{ID: "key", MType: "gauge", Value: &gauge},
			want: response{http.StatusOK, "application/json",
				"{\"id\":\"key\",\"type\":\"gauge\",\"value\":1.2345}\n"},
		},
		{
			name:   "invalid gauge metric with body",
			url:    "/update/",
			method: http.MethodPost,
			body:   &metric.Metric{ID: "key", MType: "counter", Value: &gauge},
			want: response{http.StatusBadRequest, "",
				""},
		},
		{
			name:   "invalid counter metric with body",
			url:    "/update/",
			method: http.MethodPost,
			body:   &metric.Metric{ID: "key", MType: "gauge", Delta: &counter},
			want: response{http.StatusBadRequest, "",
				""},
		},
		{
			name:   "correct value request",
			url:    "/value/",
			method: http.MethodPost,
			body:   &metric.Metric{ID: "key", MType: "gauge"},
			want: response{http.StatusOK, "application/json",
				"{\"id\":\"key\",\"type\":\"gauge\",\"value\":1.2345}\n"},
		},
		{
			name:   "no metric in inmemstore",
			url:    "/value/",
			method: http.MethodPost,
			body:   &metric.Metric{ID: "ke", MType: "gauge"},
			want: response{http.StatusNotFound, "",
				""},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tc.method, tc.url, tc.body)
			defer resp.Body.Close()
			assert.Equal(t, tc.want.statusCode, resp.StatusCode)
			assert.Equal(t, tc.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tc.want.body, get)
		})
	}
}
