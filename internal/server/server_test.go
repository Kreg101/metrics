package server

import (
	"github.com/Kreg101/metrics/internal/server/handler"
	"testing"
)

func TestCreateNewServer(t *testing.T) {
	tt := []struct {
		name string
		want *Server
	}{
		{name: "basic", want: &Server{handler.NewMux(), ""}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//assert.Equal(t, tc.want, CreateNewServer())
		})
	}
}

func TestServer_ListenAndServe(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "too much servers for a single port"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {})
	}
}
