package main

import (
	"database/sql"
	"github.com/Kreg101/metrics/internal/server/db/client"
	"github.com/Kreg101/metrics/internal/server/handler"
	"github.com/Kreg101/metrics/internal/server/inmemstore"
	"go.uber.org/zap"
)

func repInit(log *zap.SugaredLogger) (handler.Repository, error) {
	if useDB {
		db, err := sql.Open("pgx", databaseDSN)
		if err != nil {
			log.Errorf("can't connect to data base: %e", err)
			return nil, err
		}
		return client.NewClient(db), nil
	}

	return inmemstore.NewInMemStorage(storagePath, storeInterval, restore, log)

}
