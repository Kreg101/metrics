package metric

import (
	"fmt"
	"github.com/Kreg101/metrics/internal/algorithms"
	"strings"
)

type Metrics map[string]Metric

func (m Metrics) String() string {
	list := make([]string, 0)
	for k, v := range m {
		var keyValue = k + ":"
		switch v.MType {
		case "gauge":
			keyValue += algorithms.Float2String(*v.Value)
		case "counter":
			keyValue += fmt.Sprintf("%d", *v.Delta)
		}
		list = append(list, keyValue)
	}
	return strings.Join(list, ", ")
}
