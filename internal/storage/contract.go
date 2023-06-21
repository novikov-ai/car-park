package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Client interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
