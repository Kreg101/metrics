package server

import (
	"github.com/Kreg101/metrics/internal/server/handler"
	"net/http"
)

type Server struct {
	mux  *handler.Mux
	host string
}

func CreateNewServer() *Server {
	var serv = &Server{nil, ""}
	serv.mux = handler.NewMux()
	return serv
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}
