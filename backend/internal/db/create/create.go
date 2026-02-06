package create

import (
	"context"
	"fmt"
	"log"
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
	query := `CREATE DATABASE bento_db IF NOT EXISTS;`
	tag, err := conn.Exec(context.Background(), query)
	return tag.RowsAffected(), err
}

// create table
func CreateTable(conn *pgxpool.Pool, name, schema string) (int64, error) {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (%s);`, name, schema)
	tag, err := conn.Exec(context.Background(), query)
	return tag.RowsAffected(), err
}

// create tables
func CreateTables(conn *pgxpool.Pool) error {
	tables := map[string]string{
		"menu":       "id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, restaurant_id INT NOT NULL REFERENCES restaurant(id), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()",
		"version":    "id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, menu_id INT NOT NULL REFERENCES menu(id), s3_url TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), is_active BOOLEAN NOT NULL",
		"item":       "id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, version_id INT NOT NULL REFERENCES version(id), name TEXT NOT NULL, description TEXT, price NUMERIC(10,2) NOT NULL, category TEXT, modifiers JSONB",
		"users":      "id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, uuid TEXT NOT NULL",
		"restaurant": "id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, name TEXT NOT NULL, address TEXT, phone TEXT",
	}

	for name, schema := range tables {
		fmt.Printf("Creating table %s with schema %s\n", name, schema)
		rows, err := CreateTable(conn, name, schema)
		if err != nil {
			log.Printf("Error creating table %s: %v\n", name, err)
			continue
		}
		fmt.Printf("Created %d rows\n", rows)
	}

	return nil
}
