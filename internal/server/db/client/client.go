package client

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Client struct {
	db *sql.DB
}

func NewClient(db *sql.DB) Client {
	return Client{db: db}
}

func (c Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := c.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
