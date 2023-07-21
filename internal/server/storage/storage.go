package storage

type Gauge float64
type Counter int64

type Metrics map[string]interface{}

type Storage struct {
	metrics Metrics
}

// NewStorage return pointer to Storage with initialized metrics field
func NewStorage() *Storage {
	storage := &Storage{}
	storage.metrics = Metrics{}
	return storage
}

func (s *Storage) Add(key string, value interface{}) {
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
	return s.metrics
}

// Get return an element, true if it exists in map or nil, false if it's not
func (s *Storage) Get(name string) (interface{}, bool) {
	if v, ok := s.metrics[name]; ok {
		return v, ok
	}
	return nil, false
}

// CheckType returns string, because it's easier to compare result with pattern
// in handler's functions
func (s *Storage) CheckType(name string) string {
	switch s.metrics[name].(type) {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	}
	return ""
}
