package algorithms

import (
	"github.com/Kreg101/metrics/internal/metric"
	"github.com/Kreg101/metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

// I can't test this function with multiple values
// because the order of elements in map is variable
func Test_metricsToString(t *testing.T) {
	x := int64(10)
	y := 123.4
	counter := metric.Metric{ID: "x", MType: "counter", Delta: &x}
	gauge := metric.Metric{ID: "x", MType: "gauge", Value: &y}
	tt := []struct {
		name   string
		source storage.Metrics
		want   string
	}{
		{
			name:   "single counter metric",
			source: storage.Metrics{"x": counter},
			want:   "x:10",
		},
		{
			name:   "single gauge metric",
			source: storage.Metrics{"x": gauge},
			want:   "x:123.4",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Metrics2String(tc.source))
		})
	}
}

func Test_singleMetricToString(t *testing.T) {
	x := int64(10)
	y := 123.4
	counter := metric.Metric{ID: "x", MType: "counter", Delta: &x}
	gauge := metric.Metric{ID: "x", MType: "gauge", Value: &y}
	tt := []struct {
		name   string
		source storage.Metrics
		want   string
	}{
		{
			name:   "single counter metric",
			source: storage.Metrics{"x": counter},
			want:   "x:10",
		},
		{
			name:   "single gauge metric",
			source: storage.Metrics{"x": gauge},
			want:   "x:123.4",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Metrics2String(tc.source))
		})
	}
}

func Test_Float2String(t *testing.T) {
	tt := []struct {
		name string
		args float64
		want string
	}{
		{
			name: "no trim",
			args: 1.235,
			want: "1.235",
		},
		{
			name: "trim 1 digit",
			args: 1.230,
			want: "1.23",
		},
		{
			name: "trim 2 digits",
			args: 1.200,
			want: "1.2",
		},
		{
			name: "integer",
			args: 1.00000,
			want: "1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Float2String(tc.args))
		})
	}
}
