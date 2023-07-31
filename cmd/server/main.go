package main

import (
	"github.com/Kreg101/metrics/internal/server"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/Kreg101/metrics/internal/server/storage"
	"time"
)

func main() {

	parseConfiguration()

	log := logger.Default()
	defer log.Sync()

	repository, err := storage.NewStorage(storagePath, storeInterval, fileWrite, restore)
	if err != nil {
		log.Fatalf("can't initialize storage: %e", err)
	}

	if storeInterval != 0 {
		go func(s *storage.Storage, d time.Duration) {
			for range time.Tick(d) {
				s.Write()
			}
		}(repository, time.Duration(storeInterval)*time.Second)
	}

	s := server.NewServer(repository)
	err = s.Start(endpoint)
	if err != nil {
		log.Fatalf("can't start server: %e", err)
	}

}
