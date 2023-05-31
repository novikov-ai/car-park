package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func New(conn *pgx.Conn) Client {
	return &Storage{db: conn}
}

func (c *Storage) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if len(args) == 0 {
		return c.db.Query(ctx, sql)
	}

	return c.db.Query(ctx, sql, args...)
}
