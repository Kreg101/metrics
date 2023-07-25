package server

import (
	"github.com/Kreg101/metrics/internal/server/handler"
	"net/http"
)

type Server struct {
	mux  *handler.Mux
	host string
}

func NewServer(repository handler.Repository) *Server {
	serv := &Server{nil, ""}
	serv.mux = handler.NewMux(repository)
	return serv
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux.Router())
}
