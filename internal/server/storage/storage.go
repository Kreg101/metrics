package storage

import (
	"github.com/Kreg101/metrics/internal/server/constants"
)

type Gauge float64
type Counter int64

type Metrics map[string]interface{}

type Storage struct {
	metrics *Metrics
}

func NewStorage() *Storage {
	storage := &Storage{}
	storage.metrics = &Metrics{}
	return storage
}

func (s *Storage) Add(key string, value interface{}) {
	switch v := value.(type) {
	case Gauge:
		(*s.metrics)[key] = v
	case Counter:
		if val, ok := (*s.metrics)[key]; ok {
			(*s.metrics)[key] = val.(Counter) + v
		} else {
			(*s.metrics)[key] = v
		}
	}
}

func (s *Storage) GetAll() Metrics {
	return *s.metrics
}

func (s *Storage) Get(name string) (interface{}, bool) {
	if v, ok := (*s.metrics)[name]; ok {
		return v, ok
	}
	return nil, false
}

func (s *Storage) CheckType(name string) string {
	switch (*s.metrics)[name].(type) {
	case Gauge:
		return constants.GaugeType
	case Counter:
		return constants.CounterType
	}
	return ""
}

func createStorageFromMap(m map[string]interface{}) *Storage {
	s := &Storage{}
	s.metrics = &Metrics{}
	for k, v := range m {
		(*s.metrics)[k] = v
	}
	return s
}
