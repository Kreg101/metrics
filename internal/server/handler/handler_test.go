package handler

import (
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
		expected *Mux
	}{
		{
			name:     "default constructor",
			expected: &Mux{storage: storage.NewStorage()},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mux := NewMux(storage.NewStorage())
			assert.Equal(t, tc.expected, mux)
		})
	}
}

// I can't test this function with multiple values
// because the order of elements in map is variable
func Test_metricsToString(t *testing.T) {
	tt := []struct {
		name   string
		source storage.Metrics
		want   string
	}{
		{
			name:   "single counter metric",
			source: storage.Metrics{"x": storage.Counter(1)},
			want:   "x:1",
		},
		{
			name:   "single gauge metric",
			source: storage.Metrics{"x": storage.Gauge(1.340)},
			want:   "x:1.34",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, metrics2String(tc.source))
		})
	}
}

func Test_singleMetricToString(t *testing.T) {
	tt := []struct {
		name   string
		source interface{}
		want   string
	}{
		{
			name:   "counter metric",
			source: storage.Counter(1),
			want:   "1",
		},
		{
			name:   "gauge metric",
			source: storage.Gauge(1.34),
			want:   "1.34",
		},
		{
			name:   "invalid type metric",
			source: 2,
			want:   "",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, singleMetric2String(tc.source))
		})
	}
}

func Test_float2String(t *testing.T) {
	tt := []struct {
		name string
		args float64
		want string
	}{
		{
			name: "no trim",
			args: 1.235,
			want: "1.235",
		},
		{
			name: "trim 1 digit",
			args: 1.230,
			want: "1.23",
		},
		{
			name: "trim 2 digits",
			args: 1.200,
			want: "1.2",
		},
		{
			name: "integer",
			args: 1.00000,
			want: "1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, float2String(tc.args))
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
	mux := NewMux(storage.NewStorage())
	ts := httptest.NewServer(mux.Router())
	defer ts.Close()
	type response struct {
		statusCode int
		body       string
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
			want:   response{http.StatusOK, ""},
		},
		{
			name:   "correct update counter request",
			url:    "/update/counter/x/10",
			method: http.MethodPost,
			want:   response{http.StatusOK, ""},
		},
		{
			name:   "correct update gauge request",
			url:    "/update/gauge/y/1.23",
			method: http.MethodPost,
			want:   response{http.StatusOK, ""},
		},
		{
			name:   "no metric name in update request #1",
			url:    "/update/counter//10",
			method: http.MethodPost,
			want:   response{http.StatusNotFound, ""},
		},
		{
			name:   "no metric name in update request #2",
			url:    "/update/counter",
			method: http.MethodPost,
			want:   response{http.StatusNotFound, "404 page not found\n"},
		},
		{
			name:   "invalid counter type value",
			url:    "/update/counter/x/abc",
			method: http.MethodPost,
			want:   response{http.StatusBadRequest, ""},
		},
		{
			name:   "invalid gauge type value",
			url:    "/update/counter/x/abc",
			method: http.MethodPost,
			want:   response{http.StatusBadRequest, ""},
		},
		{
			name:   "invalid metric type",
			url:    "/update/counte/x/abc",
			method: http.MethodPost,
			want:   response{http.StatusBadRequest, ""},
		},
		{
			name:   "single metric request",
			url:    "/value/counter/x",
			method: http.MethodGet,
			want:   response{statusCode: http.StatusOK, body: "10"},
		},
		{
			name:   "invalid metric type request",
			url:    "/value/cor/x",
			method: http.MethodGet,
			want:   response{statusCode: http.StatusNotFound, body: ""},
		},
		{
			name:   "no metric in storage",
			url:    "/value/counter/z",
			method: http.MethodGet,
			want:   response{statusCode: http.StatusNotFound, body: ""},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, tc.method, tc.url)
			defer resp.Body.Close()
			assert.Equal(t, tc.want.statusCode, resp.StatusCode)
			assert.Equal(t, tc.want.body, get)
		})
	}
}
