package algorithms

import (
	"fmt"
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/storage"
	"math"
	"strings"
)

func Metrics2String(m storage.Metrics) string {
	list := make([]string, 0)
	for k, v := range m {
		var keyValue = k + ":"
		switch v.MType {
		case "gauge":
			keyValue += Float2String(*v.Value)
		case "counter":
			keyValue += fmt.Sprintf("%d", *v.Delta)
		}
		list = append(list, keyValue)
	}
	return strings.Join(list, ", ")
}

func SingleMetric2String(m metric.Metric) string {
	switch m.MType {
	case "gauge":
		return Float2String(*m.Value)
	case "counter":
		return fmt.Sprintf("%d", *m.Delta)
	}
	return ""
}

func Float2String(v float64) string {
	if math.Trunc(v) == v {
		return fmt.Sprintf("%.0f", v)
	}
	return strings.TrimRight(fmt.Sprintf("%.3f", v), "0")
}
