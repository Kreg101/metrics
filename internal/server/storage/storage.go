package storage

import (
	"encoding/json"
	"fmt"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/logger"
	"os"
	"sync"
	"time"
)

type Metrics map[string]metric.Metric

type Storage struct {
	metrics           Metrics
	mutex             *sync.RWMutex
	consumer          *Consumer
	producer          *Producer
	storeInterval     time.Duration
	syncWritingToFile bool
}

func NewStorage(path string, storeInterval int, writeFile, loadFromFile bool) (*Storage, error) {
	storage := &Storage{}
	storage.metrics = Metrics{}
	storage.mutex = &sync.RWMutex{}

	if !writeFile {
		return storage, nil
	}

	//var help string
	//for i, x := range strings.Split(path, "/") {
	//	if i != len(strings.Split(path, "/"))-1 && x != "" {
	//		help += "/" + x
	//		fmt.Println(help)
	//		err := os.Mkdir(help, 0666)
	//		if err != nil {
	//			fmt.Println("haahah")
	//			return nil, err
	//		}
	//	}
	//}
	//err := os.Mkdir("tmp", 0666)
	//if err != nil {
	//	fmt.Println("rrrrr")
	//	return nil, err
	//}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("lox")
		return nil, err
	}

	storage.producer = &Producer{file, json.NewEncoder(file)}
	storage.consumer = &Consumer{file, json.NewDecoder(file)}

	if loadFromFile {
		storage.metrics, err = storage.consumer.LoadFile()
		if err != nil {
			return nil, err
		}
	}

	if storeInterval == 0 {
		storage.syncWritingToFile = true
	}

	return storage, nil
}

func (s *Storage) Add(m metric.Metric) {
	log := logger.Default()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if m.MType == "counter" {
		if v, ok := s.metrics[m.ID]; ok {
			*m.Delta += *v.Delta
		}
	}
	s.metrics[m.ID] = m
	if s.syncWritingToFile {
		err := s.producer.WriteMetric(&m)
		if err != nil {
			log.Fatalf("can't add metric %v to file: %e", &m, err)
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
func (s *Storage) Get(name string) (metric.Metric, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if v, ok := s.metrics[name]; ok {
		return v, ok
	}
	return metric.Metric{}, false
}
