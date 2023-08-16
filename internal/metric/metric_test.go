package metric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetric_String(t *testing.T) {
	x := int64(10)
	y := 123.4
	counter := Metric{ID: "x", MType: "counter", Delta: &x}
	gauge := Metric{ID: "x", MType: "gauge", Value: &y}
	tt := []struct {
		name   string
		source Metrics
		want   string
	}{
		{
			name:   "single counter metric",
			source: Metrics{"x": counter},
			want:   "x:10 ",
		},
		{
			name:   "single gauge metric",
			source: Metrics{"x": gauge},
			want:   "x:123.4 ",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.source.String())
		})
	}
}

// I can't test this function with multiple values
// because the order of elements in map is variable
func TestMetrics_String(t *testing.T) {
	x := int64(10)
	y := 123.4
	counter := Metric{ID: "x", MType: "counter", Delta: &x}
	gauge := Metric{ID: "x", MType: "gauge", Value: &y}
	tt := []struct {
		name   string
		source Metrics
		want   string
	}{
		{
			name:   "single counter metric",
			source: Metrics{"x": counter},
			want:   "x:10 ",
		},
		{
			name:   "single gauge metric",
			source: Metrics{"x": gauge},
			want:   "x:123.4 ",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.source.String())
		})
	}
}
