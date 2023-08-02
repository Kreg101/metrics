package storage

import (
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"os"
	"sync"
	"testing"
)

// TestNewStorage: need to clean or delete
func TestNewStorage(t *testing.T) {
	type params struct {
		path          string
		storeInterval int
		loadFromFile  bool
		log           *zap.SugaredLogger
	}
	tt := []struct {
		name  string
		param params
		want  *Storage
	}{
		{
			name:  "basic",
			param: params{"", 0, false, nil},
			want:  &Storage{mutex: &sync.RWMutex{}, log: logger.Default()},
		},
		{
			name:  "load from empty file and sync writing",
			param: params{"tests.json", 0, true, nil},
			want:  &Storage{mutex: &sync.RWMutex{}, log: logger.Default(), filer: Filer{}, syncWritingToFile: true},
		},
		{
			name:  "load from empty file and not sync writing",
			param: params{"tests.json", 10, true, nil},
			want:  &Storage{mutex: &sync.RWMutex{}, log: logger.Default(), filer: Filer{}, syncWritingToFile: false},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want.metrics = metric.Metrics{}
			s, err := NewStorage(tc.param.path,
				tc.param.storeInterval,
				tc.param.loadFromFile,
				tc.param.log)
			s.filer = Filer{}
			require.Nil(t, err)
			assert.Equal(t, tc.want, s)
		})
	}
	defer os.Remove("test.json")
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
			base:     &Storage{mutex: &sync.RWMutex{}, metrics: metric.Metrics{}},
			expected: &Storage{mutex: &sync.RWMutex{}, metrics: map[string]metric.Metric{"key": counter}},
		},
		{
			name:     "add gauge to empty",
			value:    gauge,
			base:     &Storage{mutex: &sync.RWMutex{}, metrics: metric.Metrics{}},
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
		source metric.Metrics
		key    string
		value  metric.Metric
		ok     bool
	}{
		{
			name:   "value in storage",
			source: metric.Metrics{"c1": counter1},
			key:    "c1",
			value:  counter1,
			ok:     true,
		},
		{
			name:   "value is not in storage",
			source: metric.Metrics{"c1": counter1},
			key:    "x",
			value:  metric.Metric{},
			ok:     false,
		},
		{
			name:   "empty storage",
			source: metric.Metrics{},
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

func TestStorage_GetAll(t *testing.T) {
	x := int64(10)
	y := 1.23
	counter := metric.Metric{ID: "c", MType: "counter", Delta: &x}
	gauge := metric.Metric{ID: "g", MType: "gauge", Value: &y}
	tt := []struct {
		name string
		s    *Storage
		want metric.Metrics
	}{
		{
			name: "empty storage",
			s:    &Storage{metrics: metric.Metrics{}, mutex: &sync.RWMutex{}},
			want: metric.Metrics{},
		},
		{
			name: "not empty storage",
			s:    &Storage{metrics: metric.Metrics{"c": counter, "g": gauge}, mutex: &sync.RWMutex{}},
			want: metric.Metrics{"c": counter, "g": gauge},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.s.GetAll())
		})
	}
}

func Test_lineCounter(t *testing.T) {
	tt := []struct {
		name     string
		fileName string
		wantErr  bool
		want     int
	}{
		{
			name:     "empty file",
			fileName: "tests.json",
			wantErr:  false,
			want:     0,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.fileName)
			require.NoError(t, err)
			got, err := lineCounter(file)
			require.True(t, (tc.wantErr == true && err != nil) || (tc.wantErr == false && err == nil))
			assert.Equal(t, tc.want, got)
		})
	}
	defer os.Remove("tests.json")
}
