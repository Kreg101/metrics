package main

import (
	"github.com/Kreg101/metrics/internal/server"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/Kreg101/metrics/internal/server/storage"
)

func main() {

	parseFlags()

	log := logger.Default()
	defer log.Sync()

	s := server.NewServer(storage.NewStorage())
	
	err := s.Start(flagEndpoint)
	if err != nil {
		panic(err)
	}

}
