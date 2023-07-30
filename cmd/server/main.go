package main

import (
	"github.com/Kreg101/metrics/internal/server"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/Kreg101/metrics/internal/server/storage"
)

func main() {

	parseConfiguration()

	log := logger.Default()
	defer log.Sync()

	fileWrite = false
	repository, err := storage.NewStorage(storagePath, storeInterval, fileWrite, restore)
	if err != nil {
		//fmt.Println("here")
		log.Fatalf("can't initialize storage: %e", err)
	}

	s := server.NewServer(repository)
	err = s.Start(endpoint)
	if err != nil {
		log.Fatalf("can't start server: %e", err)
	}

}
