package entity

import (
	"fmt"
	"github.com/Kreg101/metrics/pkg/algo"
)

// Metric - единица хранения и передачи между хэндлерами и хранилищем
type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metric) String() string {
	switch m.MType {
	case "gauge":
		return algo.Float2String(*m.Value)
	case "counter":
		return fmt.Sprintf("%d", *m.Delta)
	}
	return ""
}
