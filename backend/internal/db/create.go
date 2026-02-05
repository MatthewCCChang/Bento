package db

import (
	"context"
	"os"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateConnectionPool(connections int) (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DB_URL")
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse DATABASE_URL: %v\n", err)
		os.Exit(1)
	}
	config.MaxConns = int32(connections)

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	return conn, nil
}