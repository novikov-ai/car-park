package cmd

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func MustInitDB(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = conn.Ping(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to ping db: %v", err)
		os.Exit(1)
	}

	return conn
}
