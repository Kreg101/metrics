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
	case Counter:
		(*s.metrics)[key] = v
	case Gauge:
		if val, ok := (*s.metrics)[key]; ok {
			switch help := val.(type) {
			case Gauge:
				(*s.metrics)[key] = help + v
			}
		} else {
			(*s.metrics)[key] = v
		}
	}
}
