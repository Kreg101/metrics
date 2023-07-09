package storage

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
		if val, ok := (*s.metrics)[key]; ok { // if value of counter type is already in metrics
			(*s.metrics)[key] = val.(Counter) + v // we should update it
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

// CheckType returns string, because it's easier to compare result with pattern
// in handler's functions
func (s *Storage) CheckType(name string) string {
	switch (*s.metrics)[name].(type) {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	}
	return ""
}
