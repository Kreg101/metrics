package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {

	type params struct {
		update int
		send   int
		host   string
		key    string
	}

	tt := []struct {
		name  string
		param params
		want  *Agent
	}{
		{
			name: "basic",
			param: params{
				update: 2,
				send:   10,
				host:   "http://localhost",
				key:    "",
			},
			want: &Agent{
				updateFreq: 2 * time.Second,
				sendFreq:   10 * time.Second,
				host:       "http://localhost",
				key:        "",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			newAgent := NewAgent(tc.param.update, tc.param.send, tc.param.host, tc.param.key)
			assert.Equal(t, tc.want, newAgent)
		})
	}
}
