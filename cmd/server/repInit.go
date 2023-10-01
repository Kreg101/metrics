package main

import (
	"database/sql"
	"github.com/Kreg101/metrics/internal/server/infrastructure/db"
	"github.com/Kreg101/metrics/internal/server/infrastructure/inmemstore"
	"github.com/Kreg101/metrics/internal/server/transport"
	"go.uber.org/zap"
)

func repInit(log *zap.SugaredLogger) (transport.Repository, error) {
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
