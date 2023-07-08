package handler

import (
	"github.com/Kreg101/metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMux(t *testing.T) {
	tt := []struct {
		name     string
		expected *Mux
	}{
		{
			name: "default constructor", expected: &Mux{storage: storage.NewStorage()},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mux := NewMux()
			mux.router = nil
			assert.Equal(t, tc.expected, mux)
		})
	}
}

// нужно вернуться в эту фунцию и проверить, что в storage попадают корректные значения
//func TestMux_ServeHTTP(t *testing.T) {
//	type want struct {
//		code        int
//		contentType string
//		response    string
//	}
//	tt := []struct {
//		name     string
//		method   string
//		request  string
//		expected want
//	}{
//		{name: "not POST request", method: http.MethodGet, request: "/update", expected: want{code: http.StatusMethodNotAllowed}},
//		{name: "BadRequest #1", method: http.MethodPost, request: "/ha/counter/x/0", expected: want{code: http.StatusBadRequest}},
//		{name: "BadRequest #2", method: http.MethodPost, request: "/update/gau/x/0.0", expected: want{code: http.StatusBadRequest}},
//		{name: "BadRequest #3", method: http.MethodPost, request: "/update/counter/x/abc", expected: want{code: http.StatusBadRequest}},
//		{name: "NotFound #1", method: http.MethodPost, request: "/ha/counter", expected: want{code: http.StatusNotFound}},
//		{name: "NotFound #2", method: http.MethodPost, request: "/update/counter//0", expected: want{code: http.StatusNotFound}},
//		{name: "Ok #1", method: http.MethodPost, request: "/update/counter/x/0", expected: want{code: http.StatusOK, contentType: "application/json", response: ""}},
//	}
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			request := httptest.NewRequest(tc.method, tc.request, nil)
//			w := httptest.NewRecorder()
//			mux := NewMux()
//			mux.ServeHTTP(w, request)
//
//			result := w.Result()
//
//			require.Equal(t, tc.expected.code, result.StatusCode)
//			if result.StatusCode != http.StatusOK {
//				return
//			}
//			assert.Equal(t, tc.expected.contentType, result.Header.Get("Content-Type"))
//
//			body, err := io.ReadAll(result.Body)
//			require.NoError(t, err)
//			err = result.Body.Close()
//			require.NoError(t, err)
//
//			assert.Equal(t, tc.expected.response, string(body))
//		})
//	}
//}

func Test_metricsToString(t *testing.T) {
	tt := []struct {
		name   string
		source storage.Metrics
		want   string
	}{
		{name: "single counter metric", source: storage.Metrics{"x": storage.Counter(1)}, want: "x:1"},
		{name: "single gauge metric", source: storage.Metrics{"x": storage.Gauge(1.34)}, want: "x:1.340"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, metricsToString(tc.source))
		})
	}
}

func Test_singleMetricToString(t *testing.T) {
	tt := []struct {
		name   string
		source interface{}
		want   string
	}{
		{name: "counter metric", source: storage.Counter(1), want: "1"},
		{name: "gauge metric", source: storage.Gauge(1.34), want: "1.340"},
		{name: "invalid type metric", source: 2, want: ""},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, singleMetricToString(tc.source))
		})
	}
}
