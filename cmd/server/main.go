package main

import (
	"github.com/Kreg101/metrics/internal/server"
	"github.com/Kreg101/metrics/internal/server/logger"
)

func main() {

	parseConfiguration()

	log := logger.Default()
	defer log.Sync()

	repository, err := repInit(log)
	if err != nil {
		panic(err)
	}

	s := server.NewServer(repository, log)
	err = s.Start(endpoint)
	if err != nil {
		panic(err)
	}

}


