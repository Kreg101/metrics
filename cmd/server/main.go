package main

import (
	"github.com/Kreg101/metrics/internal/server"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/Kreg101/metrics/internal/server/storage"
)

func main() {

	parseFlags()

	log := logger.New()
	defer log.Sync()

	repository := storage.NewStorage()

	s := server.NewServer(repository)
	err := s.Start(flagEndpoint)
	if err != nil {
		log.Fatalf("can't start server: %e", err)
	}

}
