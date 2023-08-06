package main

import (
	"database/sql"
	"github.com/Kreg101/metrics/internal/server"
	"github.com/Kreg101/metrics/internal/server/db/client"
	"github.com/Kreg101/metrics/internal/server/logger"
	"github.com/Kreg101/metrics/internal/server/storage"
	"time"
)

func main() {

	parseConfiguration()

	log := logger.Default()
	defer log.Sync()

	repository, err := storage.NewStorage(storagePath, storeInterval, restore, log)
	if err != nil {
		log.Fatalf("can't initialize storage: %e", err)
	}

	// Проверяем, нужно ли нам с заданном переодичностью писать данные хранилища в файл
	// если storeInterval == 0, то мы должны синхронно записывать данные в файл
	if storeInterval != 0 {
		go func(s *storage.Storage, d time.Duration) {
			for range time.Tick(d) {
				s.Write()
			}
		}(repository, time.Duration(storeInterval)*time.Second)
	}

	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	dbClient := client.NewClient(db)

	s := server.NewServer(repository, log, dbClient)
	err = s.Start(endpoint)
	if err != nil {
		panic(err)
	}

}
