package agent

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {

	tt := []struct {
		name string
		want *Agent
	}{
		{name: "basic", want: &Agent{updateFreq: 2 * time.Second, sendFreq: 10 * time.Second, host: "http://localhost", client: &resty.Client{}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			newAgent := NewAgent(2, 10, "http://localhost")
			newAgent.client = &resty.Client{}
			assert.Equal(t, tc.want, newAgent)
		})
	}
}
