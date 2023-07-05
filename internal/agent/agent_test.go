package agent

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {

	tt := []struct {
		name string
		want *Agent
	}{
		{name: "basic", want: &Agent{updateFreq: 2 * time.Second, sendFreq: 10 * time.Second, host: "http://localhost", client: http.Client{}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, NewAgent(2*time.Second, 10*time.Second, "http://localhost"))
		})
	}
}

func TestAgent_Start(t *testing.T) {
	type fields struct {
		updateFreq time.Duration
		sendFreq   time.Duration
		host       string
		stats      runtime.MemStats
		client     http.Client
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				updateFreq: tt.fields.updateFreq,
				sendFreq:   tt.fields.sendFreq,
				host:       tt.fields.host,
				stats:      tt.fields.stats,
				client:     tt.fields.client,
			}
			a.Start()
		})
	}
}

func Test_getMapOfStats(t *testing.T) {
	type args struct {
		stats *runtime.MemStats
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMapOfStats(tt.args.stats); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMapOfStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
