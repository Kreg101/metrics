package main

import (
	"database/sql"
	"github.com/Kreg101/metrics/internal/server/db/client"
	"github.com/Kreg101/metrics/internal/server/handler"
	"github.com/Kreg101/metrics/internal/server/inMemStore"
	"go.uber.org/zap"
)

func repInit(repository handler.Repository, log *zap.SugaredLogger) error {
	var err error
	if useDB {
		db, err := sql.Open("pgx", databaseDSN)
		if err != nil {
			log.Errorf("can't connect to data base: %e", err)
			return err
		}
		repository = client.NewClient(db)
		return nil
	}

	repository, err = inMemStore.NewInMemStorage(storagePath, storeInterval, restore, log)
	if err != nil {
		return err
	}

	return nil
}
