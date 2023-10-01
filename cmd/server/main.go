package main

import (
	"github.com/Kreg101/metrics/internal/server/transport"
	"github.com/Kreg101/metrics/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

func main() {

	parseConfiguration()

	log := logger.Default()
	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			log.Fatalf("can't initialize logger: %v", err)
		}
	}(log)

	repository, err := repInit(log)
	if err != nil {
		panic(err)
	}

	t := transport.New(repository, log, key)
	err = http.ListenAndServe(endpoint, t.Router())
	if err != nil {
		panic(err)
	}
}
