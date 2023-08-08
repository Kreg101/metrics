package client

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type DBClient struct {
	db *sql.DB
}

// Open подключается к базе данных с заданным dsn и возвращает клиента, подлкюченного к
// этой базе данных
func Open(dsn string) (DBClient, error) {
	c := DBClient{}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return c, err
	}

	// Можно еще дополнительно пинговать для проверки соединения

	c.db = db
	return c, nil
}

// Close закрывает конекшн к базе данных
func (c DBClient) Close() {
	c.db.Close()
}

//func NewClient(db *sql.DB) Client {
//	return Client{db: db}
//}

func (c DBClient) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := c.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
