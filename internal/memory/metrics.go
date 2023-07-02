package memory

type Gauge float64
type Counter int64

type Metrics map[string]interface{}

func NewMetrics() *Metrics {
	m := make(Metrics, 0)
	return &m
}

func (m *Metrics) AddMetric(key string, x interface{}) {
	switch v := x.(type) {
	case Counter:
		(*m)[key] = v
	case Gauge:
		(*m)[key] = v
	}
}
