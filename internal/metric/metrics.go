package metric

import (
	"fmt"
	"github.com/Kreg101/metrics/internal/algorithms"
	"strings"
)

type Metrics map[string]Metric

func (m Metrics) String() string {
	sb := strings.Builder{}
	for k, v := range m {
		sb.WriteString(k)
		sb.WriteString(":")
		switch v.MType {
		case "gauge":
			_, _ = sb.WriteString(algorithms.Float2String(*v.Value))
		case "counter":
			_, _ = sb.WriteString(fmt.Sprintf("%d", *v.Delta))
		}
	}
	return sb.String()
}
