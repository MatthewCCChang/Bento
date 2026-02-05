package create

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Create connection pool to the database
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

// create database
func CreateDatabase(conn *pgxpool.Pool) (int64, error) {
	query := `CREATE DATABASE bento_db;`
	tag, err := conn.Exec(context.Background(), query)
	return tag.RowsAffected(), err
}

// create table
func CreateTable(conn *pgxpool.Pool, name, schema string) (int64, error) {
	query := fmt.Sprintf(`CREATRE TABLE IF NOT EXISTS %s.%s;`, name, schema)
	tag, err := conn.Exec(context.Background(), query)
	return tag.RowsAffected(), err
}

// create tables
func CreateTables(conn *pgxpool.Pool) error {
	tables := map[string]string{
		"menu":       "id SERIAL PRIMARY KEY, version_id TEXT NOT NULL, restaurant_id TEXT NOT NULL, updated_at TIMESTAMPTZ NOT NULL",
		"version":    "id SERIAL PRIMARY KEY, menu_id NUMERIC NOT NULL, s3_url TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL, is_active BOOLEAN NOT NULL",
		"item":       "id SERIAL PRIMARY KEY, version_id NUMERIC NOT NULL, name TEXT NOT NULL, description TEXT, price NUMERIC NOT NULL, category TEXT, modifiers JSONB",
		"user":       "id SERIAL PRIMARY KEY, uuid TEXT NOT NULL",
		"restaurant": "id SERIAL PRIMARY KEY, name TEXT NOT NULL, address TEXT, phone TEXT",
	}

	for name, schema := range tables {
		_, err := CreateTable(conn, name, schema)
		if err != nil {
			return err
		}
	}

	return nil
}
