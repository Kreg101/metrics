package inmemstore

import (
	"context"
	"encoding/json"
	"github.com/Kreg101/metrics/internal/server/entity"
	"github.com/Kreg101/metrics/pkg/logger"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

// InMemStorage имплементирует интерфейс Repository, поэтому может быть использован
// как хранилище метрик
type InMemStorage struct {
	metrics           entity.Metrics
	mutex             *sync.RWMutex
	log               *zap.SugaredLogger
	filer             Filer
	storeInterval     time.Duration
	syncWritingToFile bool
}

// NewInMemStorage returns initialized inmemstore pointer
func NewInMemStorage(path string, storeInterval int,
	loadFromFile bool, log *zap.SugaredLogger) (*InMemStorage, error) {

	storage := &InMemStorage{}
	storage.metrics = entity.Metrics{}
	storage.mutex = &sync.RWMutex{}

	// инициализируем логер, если он не передан, то используем дефолтный
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

	storage.filer = Filer{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}

	// проверяем, нужно ли нам загружать данные из файла при старте
	if loadFromFile {
		storage.metrics, err = storage.filer.load()
		if err != nil {
			return nil, err
		}
	}

	// проверяем, нужно ли синхронно писать в файл или делать это с заданной периодичностью
	if storeInterval == 0 {
		storage.syncWritingToFile = true
	} else {
		go func(s *InMemStorage, d time.Duration) {
			for range time.Tick(d) {
				s.Write()
			}
		}(storage, time.Duration(storeInterval)*time.Second)
	}

	return storage, nil
}

// Add (add entity to inmemstore)
func (s *InMemStorage) Add(ctx context.Context, m entity.Metric) {
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
			s.log.Errorf("can't add entity %v to file: %v", m, err)
		}
	}
}

// Get return an element, true if it exists in map or nil, false if it's not
func (s *InMemStorage) Get(ctx context.Context, name string) (entity.Metric, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if v, ok := s.metrics[name]; ok {
		return v, ok
	}
	return entity.Metric{}, false
}

// GetAll returns all metrics from inmemstore
func (s *InMemStorage) GetAll(ctx context.Context) entity.Metrics {
	s.mutex.RLock()
	duplicate := make(entity.Metrics, len(s.metrics))
	for k, v := range s.metrics {
		duplicate[k] = v
	}
	s.mutex.RUnlock()
	return duplicate
}

// Ping for in-memory inmemstore is default true
// because it doesn't need a connection. I use this function
// for common interface
func (s *InMemStorage) Ping(ctx context.Context) error {
	return nil
}
