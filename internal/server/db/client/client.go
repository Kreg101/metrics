package client

import (
	"context"
	"database/sql"
	"time"
)

type Client struct {
	db *sql.DB
}

func NewClient(init string) Client {
	var err error
	c := Client{}
	if init == "" {
		return c
	}

	c.db, err = sql.Open("pgx", init)
	if err != nil {
		panic(err)
	}
	defer c.db.Close()

	return c
}

func (c Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := c.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
