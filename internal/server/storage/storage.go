package storage

import (
	"encoding/json"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/logger"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

type Storage struct {
	metrics           metric.Metrics
	mutex             *sync.RWMutex
	log               *zap.SugaredLogger
	filer             Filer
	storeInterval     time.Duration
	syncWritingToFile bool
}

func NewStorage(path string, storeInterval int, loadFromFile bool, log *zap.SugaredLogger) (*Storage, error) {
	storage := &Storage{}
	storage.metrics = metric.Metrics{}
	storage.mutex = &sync.RWMutex{}

	if log == nil {
		storage.log = logger.Default()
	} else {
		storage.log = log
	}

	// проверяем, нужно ли нам работать с файлом
	if path == "" {
		return storage, nil
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	storage.filer = Filer{file, json.NewEncoder(file), json.NewDecoder(file)}

	// проверяем, нужно ли нам загружать данные из файла при старте
	if loadFromFile {
		storage.metrics, err = storage.filer.load()
		if err != nil {
			return nil, err
		}
	}

	// проверяем, нжуно ли синхронно писать в файл
	if storeInterval == 0 {
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
	if s.syncWritingToFile {
		err := s.filer.writeMetric(&m)
		if err != nil {
			s.log.Errorf("can't add metric %v to file: %e", &m, err)
		}
	}
}

func (s *Storage) GetAll() metric.Metrics {
	s.mutex.RLock()
	duplicate := make(metric.Metrics, len(s.metrics))
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
