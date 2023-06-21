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

func (s *Storage) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if len(args) == 0 {
		return s.db.Query(ctx, sql)
	}

	return s.db.Query(ctx, sql, args...)
}

func (s *Storage) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if len(args) == 0 {
		s.db.QueryRow(ctx, sql)
	}

	return s.db.QueryRow(ctx, sql, args...)
}
