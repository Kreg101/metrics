package server

import (
	"github.com/Kreg101/metrics/internal/server/handler"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	mux *handler.Mux
}

func NewServer(repository handler.Repository, log *zap.SugaredLogger, key string) *Server {
	return &Server{handler.NewMux(repository, log, key)}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux.Router())
}
