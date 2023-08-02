package server

import (
	"github.com/Kreg101/metrics/internal/server/handler"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	mux  *handler.Mux
	host string
}

func NewServer(repository handler.Repository, log *zap.SugaredLogger) *Server {
	serv := &Server{nil, ""}
	serv.mux = handler.NewMux(repository, log)
	return serv
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux.Router())
}
