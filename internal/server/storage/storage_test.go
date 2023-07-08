package storage

import (
	"github.com/Kreg101/metrics/internal/server/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewStorage(t *testing.T) {
	tt := []struct {
		name string
		want *Storage
	}{
		{name: "basic", want: &Storage{}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want.metrics = &Metrics{}
			assert.Equal(t, tc.want, NewStorage())
		})
	}
}

func Test_createStorageFromMap(t *testing.T) {
	tt := []struct {
		name   string
		source map[string]interface{}
		want   *Storage
	}{
		{name: "single argument", source: map[string]interface{}{"key": Counter(1)}, want: NewStorage()},
		{name: "more arguments", source: map[string]interface{}{"x": Counter(1), "y": Gauge(2.34), "z": Counter(1)}, want: NewStorage()},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.source {
				tc.want.Add(k, v)
			}
			assert.Equal(t, tc.want, createStorageFromMap(tc.source))
		})
	}
}

func TestStorage_Add(t *testing.T) {

	tt := []struct {
		name     string
		key      string
		value    interface{}
		base     *Storage
		expected *Storage
	}{
		{name: "add counter to empty", key: "key", value: Counter(1),
			base: NewStorage(), expected: createStorageFromMap(map[string]interface{}{"key": Counter(1)})},
		{name: "add gauge to empty", key: "key", value: Gauge(1.0),
			base: NewStorage(), expected: createStorageFromMap(map[string]interface{}{"key": Gauge(1.0)})},
		{name: "add counter to something", key: "key", value: Counter(-2),
			base:     createStorageFromMap(map[string]interface{}{"key": Counter(1)}),
			expected: createStorageFromMap(map[string]interface{}{"key": Counter(-1)})},
		{name: "add gauge to something", key: "key", value: Gauge(3.0),
			base:     createStorageFromMap(map[string]interface{}{"key": Gauge(-2.0)}),
			expected: createStorageFromMap(map[string]interface{}{"key": Gauge(3.0)})},
		{name: "a lot of things", key: "key", value: Gauge(1.1),
			base:     createStorageFromMap(map[string]interface{}{"x": Counter(1), "y": Gauge(3.13), "z": Counter(-1)}),
			expected: createStorageFromMap(map[string]interface{}{"x": Counter(1), "y": Gauge(3.13), "z": Counter(-1), "key": Gauge(1.1)})},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			tc.base.Add(tc.key, tc.value)
			assert.Equal(t, tc.expected, tc.base)
		})
	}
}

func TestStorage_GetAll(t *testing.T) {

	tt := []struct {
		name string
		base map[string]interface{}
		want Metrics
	}{
		{name: "basis test #1", base: map[string]interface{}{"key": Counter(1)}, want: Metrics{"key": Counter(1)}},
		{name: "basis test #2", base: map[string]interface{}{"key": Gauge(1.23)}, want: Metrics{"key": Gauge(1.23)}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := NewStorage()
			for k, v := range tc.base {
				s.Add(k, v)
			}
			assert.Equal(t, tc.want, s.GetAll())
		})
	}
}

func TestStorage_Get(t *testing.T) {
	tt := []struct {
		name   string
		source map[string]interface{}
		key    string
		value  interface{}
		ok     bool
	}{
		{name: "value in storage", source: map[string]interface{}{"x": Counter(1), "y": Gauge(2.0)}, key: "x", value: Counter(1), ok: true},
		{name: "value is not in storage", source: map[string]interface{}{}, key: "x", value: nil, ok: false},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := createStorageFromMap(tc.source)
			res, ok := s.Get(tc.key)
			require.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.value, res)
		})
	}
}

func TestStorage_CheckType(t *testing.T) {
	tt := []struct {
		name   string
		source map[string]interface{}
		key    string
		want   string
	}{
		{name: "counter value in storage", source: map[string]interface{}{"x": Counter(1), "y": Gauge(2.0)}, key: "x", want: constants.CounterType},
		{name: "gauge value in storage", source: map[string]interface{}{"x": Counter(1), "y": Gauge(2.0)}, key: "y", want: constants.GaugeType},
		{name: "value is not in storage", source: map[string]interface{}{}, key: "x", want: ""},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := createStorageFromMap(tc.source)
			assert.Equal(t, tc.want, s.CheckType(tc.key))
		})
	}
}
