package storage

import (
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestNewStorage(t *testing.T) {
	type params struct {
		path          string
		storeInterval int
		writeFile     bool
		loadFromFile  bool
	}
	tt := []struct {
		name  string
		param params
		want  *Storage
	}{
		{
			name:  "basic",
			param: params{"", 0, false, false},
			want:  &Storage{mutex: &sync.RWMutex{}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want.metrics = Metrics{}
			s, err := NewStorage(tc.param.path, tc.param.storeInterval, tc.param.writeFile, tc.param.loadFromFile)
			require.Nil(t, err)
			assert.Equal(t, tc.want, s)
		})
	}
}

func TestStorage_Add(t *testing.T) {
	x := int64(10)
	z := int64(20)
	y := 123.4
	counter := metric.Metric{ID: "key", MType: "counter", Delta: &x}
	gauge := metric.Metric{ID: "key", MType: "gauge", Value: &y}
	n := metric.Metric{ID: "new", MType: "counter", Delta: &x}
	result1 := metric.Metric{ID: "key", MType: "counter", Delta: &z}
	result2 := metric.Metric{ID: "key", MType: "gauge", Value: &y}
	tt := []struct {
		name     string
		value    metric.Metric
		base     *Storage
		expected *Storage
	}{
		{
			name:     "add counter to empty",
			value:    counter,
			base:     &Storage{mutex: &sync.RWMutex{}, metrics: Metrics{}},
			expected: &Storage{mutex: &sync.RWMutex{}, metrics: map[string]metric.Metric{"key": counter}},
		},
		{
			name:     "add gauge to empty",
			value:    gauge,
			base:     &Storage{mutex: &sync.RWMutex{}, metrics: Metrics{}},
			expected: &Storage{mutex: &sync.RWMutex{}, metrics: map[string]metric.Metric{"key": gauge}},
		},
		{
			name:     "add counter to something",
			value:    counter,
			base:     &Storage{mutex: &sync.RWMutex{}, metrics: map[string]metric.Metric{"key": counter}},
			expected: &Storage{mutex: &sync.RWMutex{}, metrics: map[string]metric.Metric{"key": result1}},
		},
		{
			name:     "add counter to something",
			value:    gauge,
			base:     &Storage{mutex: &sync.RWMutex{}, metrics: map[string]metric.Metric{"key": gauge}},
			expected: &Storage{mutex: &sync.RWMutex{}, metrics: map[string]metric.Metric{"key": result2}},
		},
		{
			name:     "add new",
			value:    n,
			base:     &Storage{mutex: &sync.RWMutex{}, metrics: map[string]metric.Metric{"key": gauge}},
			expected: &Storage{mutex: &sync.RWMutex{}, metrics: map[string]metric.Metric{"key": gauge, "new": n}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.base.Add(tc.value)
			assert.Equal(t, tc.expected, tc.base)
		})
	}
}

func TestStorage_Get(t *testing.T) {
	x := int64(10)
	counter1 := metric.Metric{ID: "c1", MType: "counter", Delta: &x}

	tt := []struct {
		name   string
		source Metrics
		key    string
		value  metric.Metric
		ok     bool
	}{
		{
			name:   "value in storage",
			source: Metrics{"c1": counter1},
			key:    "c1",
			value:  counter1,
			ok:     true,
		},
		{
			name:   "value is not in storage",
			source: Metrics{"c1": counter1},
			key:    "x",
			value:  metric.Metric{},
			ok:     false,
		},
		{
			name:   "empty storage",
			source: Metrics{},
			key:    "x",
			value:  metric.Metric{},
			ok:     false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := Storage{mutex: &sync.RWMutex{}, metrics: tc.source}
			res, ok := s.Get(tc.key)
			require.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.value, res)
		})
	}
}
