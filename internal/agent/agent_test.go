package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {

	tt := []struct {
		name string
		want *Agent
	}{
		{
			name: "basic",
			want: &Agent{
				updateFreq: 2 * time.Second,
				sendFreq:   10 * time.Second,
				host:       "http://localhost",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			newAgent := NewAgent(2, 10, "http://localhost")
			assert.Equal(t, tc.want, newAgent)
		})
	}
}
