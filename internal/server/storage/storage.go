package storage

import (
	"github.com/Kreg101/metrics/internal/metric"
	"os"
	"sync"
	"time"
)

type Metrics map[string]metric.Metric

type Storage struct {
	metrics           Metrics
	mutex             *sync.RWMutex
	file              *os.File
	storeInterval     time.Duration
	syncWritingToFile bool
}

// NewStorage return pointer to Storage with initialized metrics field
//func NewStorage() *Storage {
//	storage := &Storage{}
//	storage.metrics = Metrics{}
//	storage.mutex = &sync.RWMutex{}
//	storage.syncWritingToFile = false
//	storage.file = nil
//	return storage
//}

func NewStorage(path string, storeInterval int, writeFile, loadFromFile bool) (*Storage, error) {
	storage := &Storage{}
	storage.metrics = Metrics{}
	storage.mutex = &sync.RWMutex{}

	if !writeFile {
		return storage, nil
	}

	var err error
	storage.file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	if loadFromFile {
		// TODO load previous metrics from file
	}

	if storeInterval != 0 {
		// TODO write special goroutine to read from file
	} else {
		storage.syncWritingToFile = true
	}

	return storage, nil
}

func (s *Storage) Add(m metric.Metric) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if m.MType == "counter" {
		if v, ok := s.metrics[m.ID]; ok {
			*m.Delta += *v.Delta
		}
	}
	s.metrics[m.ID] = m
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
func (s *Storage) Get(name string) (metric.Metric, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if v, ok := s.metrics[name]; ok {
		return v, ok
	}
	return metric.Metric{}, false
}

// CheckType returns string, because it's easier to compare result with pattern
// in handler's functions
//func (s *Storage) CheckType(name string) string {
//	s.mutex.RLock()
//	defer s.mutex.RUnlock()
//	switch s.metrics[name].(type) {
//	case Gauge:
//		return "gauge"
//	case Counter:
//		return "counter"
//	}
//	return ""
//}
