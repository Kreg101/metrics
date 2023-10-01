package inmemstore

import (
	"context"
	"github.com/Kreg101/metrics/internal/server/entity"
	"github.com/Kreg101/metrics/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"os"
	"sync"
	"testing"
)

// TestNewInMemStorage: need to clean or delete
func TestNewInMemStorage(t *testing.T) {
	type params struct {
		path          string
		storeInterval int
		loadFromFile  bool
		log           *zap.SugaredLogger
	}
	tt := []struct {
		name  string
		param params
		want  *InMemStorage
	}{
		{
			name:  "basic",
			param: params{"", 0, false, nil},
			want:  &InMemStorage{mutex: &sync.RWMutex{}, log: logger.Default()},
		},
		{
			name:  "load from empty file and sync writing",
			param: params{"tests.json", 0, true, nil},
			want:  &InMemStorage{mutex: &sync.RWMutex{}, log: logger.Default(), filer: Filer{}, syncWritingToFile: true},
		},
		{
			name:  "load from empty file and not sync writing",
			param: params{"tests.json", 10, true, nil},
			want:  &InMemStorage{mutex: &sync.RWMutex{}, log: logger.Default(), filer: Filer{}, syncWritingToFile: false},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want.metrics = entity.Metrics{}
			s, err := NewInMemStorage(tc.param.path,
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

func TestInMemStorage_Add(t *testing.T) {
	x := int64(10)
	z := int64(20)
	y := 123.4
	counter := entity.Metric{ID: "key", MType: "counter", Delta: &x}
	gauge := entity.Metric{ID: "key", MType: "gauge", Value: &y}
	n := entity.Metric{ID: "new", MType: "counter", Delta: &x}
	result1 := entity.Metric{ID: "key", MType: "counter", Delta: &z}
	result2 := entity.Metric{ID: "key", MType: "gauge", Value: &y}
	tt := []struct {
		name     string
		value    entity.Metric
		base     *InMemStorage
		expected *InMemStorage
	}{
		{
			name:     "add counter to empty",
			value:    counter,
			base:     &InMemStorage{mutex: &sync.RWMutex{}, metrics: entity.Metrics{}},
			expected: &InMemStorage{mutex: &sync.RWMutex{}, metrics: map[string]entity.Metric{"key": counter}},
		},
		{
			name:     "add gauge to empty",
			value:    gauge,
			base:     &InMemStorage{mutex: &sync.RWMutex{}, metrics: entity.Metrics{}},
			expected: &InMemStorage{mutex: &sync.RWMutex{}, metrics: map[string]entity.Metric{"key": gauge}},
		},
		{
			name:     "add counter to something",
			value:    counter,
			base:     &InMemStorage{mutex: &sync.RWMutex{}, metrics: map[string]entity.Metric{"key": counter}},
			expected: &InMemStorage{mutex: &sync.RWMutex{}, metrics: map[string]entity.Metric{"key": result1}},
		},
		{
			name:     "add counter to something",
			value:    gauge,
			base:     &InMemStorage{mutex: &sync.RWMutex{}, metrics: map[string]entity.Metric{"key": gauge}},
			expected: &InMemStorage{mutex: &sync.RWMutex{}, metrics: map[string]entity.Metric{"key": result2}},
		},
		{
			name:     "add new",
			value:    n,
			base:     &InMemStorage{mutex: &sync.RWMutex{}, metrics: map[string]entity.Metric{"key": gauge}},
			expected: &InMemStorage{mutex: &sync.RWMutex{}, metrics: map[string]entity.Metric{"key": gauge, "new": n}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.base.Add(context.Background(), tc.value)
			assert.Equal(t, tc.expected, tc.base)
		})
	}
}

func TestInMemStorage_Get(t *testing.T) {
	x := int64(10)
	counter1 := entity.Metric{ID: "c1", MType: "counter", Delta: &x}

	tt := []struct {
		name   string
		source entity.Metrics
		key    string
		value  entity.Metric
		ok     bool
	}{
		{
			name:   "value in inmemstore",
			source: entity.Metrics{"c1": counter1},
			key:    "c1",
			value:  counter1,
			ok:     true,
		},
		{
			name:   "value is not in inmemstore",
			source: entity.Metrics{"c1": counter1},
			key:    "x",
			value:  entity.Metric{},
			ok:     false,
		},
		{
			name:   "empty inmemstore",
			source: entity.Metrics{},
			key:    "x",
			value:  entity.Metric{},
			ok:     false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := InMemStorage{mutex: &sync.RWMutex{}, metrics: tc.source}
			res, ok := s.Get(context.Background(), tc.key)
			require.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.value, res)
		})
	}
}

func TestInMemStorage_GetAll(t *testing.T) {
	x := int64(10)
	y := 1.23
	counter := entity.Metric{ID: "c", MType: "counter", Delta: &x}
	gauge := entity.Metric{ID: "g", MType: "gauge", Value: &y}
	tt := []struct {
		name string
		s    *InMemStorage
		want entity.Metrics
	}{
		{
			name: "empty inmemstore",
			s:    &InMemStorage{metrics: entity.Metrics{}, mutex: &sync.RWMutex{}},
			want: entity.Metrics{},
		},
		{
			name: "not empty inmemstore",
			s:    &InMemStorage{metrics: entity.Metrics{"c": counter, "g": gauge}, mutex: &sync.RWMutex{}},
			want: entity.Metrics{"c": counter, "g": gauge},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.s.GetAll(context.Background()))
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
