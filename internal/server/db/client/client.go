package client

import (
	"context"
	"database/sql"
	"github.com/Kreg101/metrics/internal/metric"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type DBClient struct {
	db *sql.DB
}

func NewClient(db *sql.DB) DBClient {
	return DBClient{db: db}
}

func (c DBClient) Add(metric metric.Metric) {
	//TODO implement me
	panic("implement me")
}

func (c DBClient) GetAll() metric.Metrics {
	//TODO implement me
	panic("implement me")
}

func (c DBClient) Get(name string) (metric.Metric, bool) {
	//TODO implement me
	panic("implement me")
}

// Open подключается к базе данных с заданным dsn и возвращает клиента, подлкюченного к
// этой базе данных
//func Open(dsn string) (DBClient, error) {
//	c := DBClient{}
//
//	db, err := sql.Open("pgx", dsn)
//	if err != nil {
//		return c, err
//	}
//
//	// Можно еще дополнительно пинговать для проверки соединения
//
//	c.db = db
//	return c, nil
//}

// Close закрывает соединение с базой данных
func (c DBClient) Close() {
	c.db.Close()
}

// Ping проверяет соединение с базой данных
func (c DBClient) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := c.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
