package storage

import (
	"fmt"
	"github.com/Kreg101/metrics/internal/server/constants"
	"strings"
)

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
	//fmt.Println(key, " ", value)
	switch v := value.(type) {
	case Counter:
		fmt.Println("i'm here")
		(*s.metrics)[key] = v
	case Gauge:
		if val, ok := (*s.metrics)[key]; ok {
			(*s.metrics)[key] = val.(Gauge) + v
		} else {
			(*s.metrics)[key] = v
		}
	}
}

// GetAll I should read about string.Builder
func (s *Storage) GetAllString() string {
	fmt.Println(len(*s.metrics))
	list := make([]string, 0)
	for k, v := range *s.metrics {
		var keyValue = k + ":"
		switch res := v.(type) {
		case Gauge:
			keyValue += fmt.Sprintf("%f", res)
		case Counter:
			keyValue += fmt.Sprintf("%d", res)
		}
		list = append(list, keyValue)
	}
	return strings.Join(list, ", ")
}

func (s *Storage) GetString(name string) (string, bool) {
	if v, ok := (*s.metrics)[name]; ok {
		switch x := v.(type) {
		case Gauge:
			return fmt.Sprintf("%f", x), ok
		case Counter:
			return fmt.Sprintf("%d", x), ok
		}
	}
	return "", false
}

func (s *Storage) CheckType(name string) string {
	switch (*s.metrics)[name].(type) {
	case Gauge:
		return constants.GaugeType
	case Counter:
		return constants.CounterType
	}
	return ""
}
