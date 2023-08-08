package main

import (
	"github.com/Kreg101/metrics/internal/server"
	"github.com/Kreg101/metrics/internal/server/handler"
	"github.com/Kreg101/metrics/internal/server/logger"
)

func main() {

	parseConfiguration()

	log := logger.Default()
	defer log.Sync()

	var repository handler.Repository

	err := repInit(repository, log)
	if err != nil {
		panic(err)
	}

	s := server.NewServer(repository, log)
	err = s.Start(endpoint)
	if err != nil {
		panic(err)
	}
}
