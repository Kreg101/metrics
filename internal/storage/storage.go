package storage

type Storage struct {
	metrics *Metrics
}

func NewStorage() *Storage {
	storage := &Storage{}
	storage.metrics = NewMetrics()
	return storage
}

func (storage *Storage) Add(key string, value interface{}) {
	storage.metrics.AddMetric(key, value)
}
