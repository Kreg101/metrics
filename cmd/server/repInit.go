package main

import (
	"database/sql"
	"github.com/Kreg101/metrics/internal/server/db"
	"github.com/Kreg101/metrics/internal/server/handler"
	"github.com/Kreg101/metrics/internal/server/inmemstore"
	"go.uber.org/zap"
)

func repInit(log *zap.SugaredLogger) (handler.Repository, error) {
	if useDB {
		conn, err := sql.Open("pgx", databaseDSN)
		if err != nil {
			log.Errorf("can't connect to data base: %e", err)
			return nil, err
		}
		return db.NewStorage(conn, log)
	}

	return inmemstore.NewInMemStorage(storagePath, storeInterval, restore, log)
}
