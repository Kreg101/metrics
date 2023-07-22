package storage

import "sync"

type Gauge float64
type Counter int64

type Metrics map[string]interface{}

type Storage struct {
	metrics Metrics
	mutex   *sync.RWMutex
}

// NewStorage return pointer to Storage with initialized metrics field
func NewStorage() *Storage {
	storage := &Storage{}
	storage.metrics = Metrics{}
	storage.mutex = &sync.RWMutex{}
	return storage
}

func (s *Storage) Add(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	switch v := value.(type) {
	case Gauge:
		s.metrics[key] = v
	case Counter:
		// if value of counter type is already in metrics - update it
		if val, ok := s.metrics[key]; ok {
			s.metrics[key] = val.(Counter) + v
		} else {
			s.metrics[key] = v
		}
	}
}

func (s *Storage) GetAll() Metrics {
	s.mutex.RLock()
	duplicate := make(Metrics, len(s.metrics))
	for k, v := range s.metrics {
		duplicate[k] = v
	}
	s.mutex.RUnlock()
	return duplicate
}

// Get return an element, true if it exists in map or nil, false if it's not
func (s *Storage) Get(name string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if v, ok := s.metrics[name]; ok {
		return v, ok
	}
	return nil, false
}

// CheckType returns string, because it's easier to compare result with pattern
// in handler's functions
func (s *Storage) CheckType(name string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	switch s.metrics[name].(type) {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	}
	return ""
}
