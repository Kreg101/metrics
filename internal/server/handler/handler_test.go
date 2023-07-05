package handler

import (
	"github.com/Kreg101/metrics/internal/server/constants"
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
			name: "basic", expected: &Mux{storage: storage.NewStorage()},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, NewMux())
		})
	}
}

// нужно вернуться в эту фунцию и проверить, что в storage попадают корректные значения
func TestMux_ServeHTTP(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    string
	}
	tt := []struct {
		name     string
		method   string
		request  string
		expected want
	}{
		{name: "not POST request", method: http.MethodGet, request: "/update", expected: want{code: http.StatusMethodNotAllowed}},
		{name: "BadRequest #1", method: http.MethodPost, request: "/ha/counter/x/0", expected: want{code: http.StatusBadRequest}},
		{name: "BadRequest #2", method: http.MethodPost, request: "/update/gau/x/0.0", expected: want{code: http.StatusBadRequest}},
		{name: "BadRequest #3", method: http.MethodPost, request: "/update/counter/x/abc", expected: want{code: http.StatusBadRequest}},
		{name: "NotFound #1", method: http.MethodPost, request: "/ha/counter", expected: want{code: http.StatusNotFound}},
		{name: "NotFound #2", method: http.MethodPost, request: "/update/counter//0", expected: want{code: http.StatusNotFound}},
		{name: "Ok #1", method: http.MethodPost, request: "/update/counter/x/0", expected: want{code: http.StatusOK, contentType: "application/json", response: ""}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, tc.request, nil)
			w := httptest.NewRecorder()
			mux := NewMux()
			mux.ServeHTTP(w, request)

			result := w.Result()

			require.Equal(t, tc.expected.code, result.StatusCode)
			if result.StatusCode != http.StatusOK {
				return
			}
			assert.Equal(t, tc.expected.contentType, result.Header.Get("Content-Type"))

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tc.expected.response, string(body))
		})
	}
}

func Test_removeFirstLastEmptyElements(t *testing.T) {
	tt := []struct {
		name     string
		data     []string
		expected []string
	}{
		{name: "nothing to do", data: []string{"abc", "bc", "d"}, expected: []string{"abc", "bc", "d"}},
		{name: "remove only first elem", data: []string{"", "a", "b"}, expected: []string{"a", "b"}},
		{name: "remove only last elem", data: []string{"a", "b", ""}, expected: []string{"a", "b"}},
		{name: "remove first and last elem", data: []string{"", "b", ""}, expected: []string{"b"}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, removeFirstLastEmptyElements(tc.data))
		})
	}
}

func Test_requestValidation(t *testing.T) {
	type answer struct {
		inter interface{}
		err   constants.Error
	}
	tt := []struct {
		name     string
		data     []string
		expected answer
	}{
		{name: "no metric name", data: []string{"update", "counter"}, expected: answer{nil, constants.NoMetricNameError}},
		{name: "no update in request", data: []string{"ha", "counter", "fun", "10"}, expected: answer{nil, constants.InvalidRequestError}},
		{name: "metric name is empty", data: []string{"update", "counter", "", "10"}, expected: answer{nil, constants.NoMetricNameError}},
		{name: "wrong length of request", data: []string{"update", "counter", "mem", "10", "x"}, expected: answer{nil, constants.InvalidRequestError}},
		{name: "invalid metric type", data: []string{"update", "counte", "ha", "0"}, expected: answer{nil, constants.InvalidMetricTypeError}},
		{name: "no metric name", data: []string{"update", "counter"}, expected: answer{nil, constants.NoMetricNameError}},
		{name: "invalid value type #1", data: []string{"update", "counter", "x", "1.00"}, expected: answer{nil, constants.InvalidValueError}},
		{name: "invalid value type #2", data: []string{"update", "gauge", "x", "abc"}, expected: answer{nil, constants.InvalidValueError}},
		{name: "correct counter request", data: []string{"update", "counter", "x", "1"}, expected: answer{storage.Counter(1), constants.NoError}},
		{name: "correct gauge request", data: []string{"update", "gauge", "x", "1.0"}, expected: answer{storage.Gauge(1.0), constants.NoError}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ans, err := requestValidation(tc.data)
			require.Equal(t, tc.expected.err, err, err)
			assert.Equal(t, tc.expected.inter, ans)
		})
	}
}
