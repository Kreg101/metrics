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
		{name: "basic", want: &Agent{updateFreq: 2 * time.Second, sendFreq: 10 * time.Second, host: "http://localhost"}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, NewAgent(2*time.Second, 10*time.Second, "http://localhost"))
		})
	}
}
