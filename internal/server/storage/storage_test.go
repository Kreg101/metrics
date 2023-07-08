package storage

import (
	"github.com/stretchr/testify/assert"
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

func createStorageFromMap(m map[string]interface{}) *Storage {
	s := &Storage{}
	s.metrics = &Metrics{}
	for k, v := range m {
		(*s.metrics)[k] = v
	}
	return s
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

func TestStorage_GetAllString(t *testing.T) {

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
