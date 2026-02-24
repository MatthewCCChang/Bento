package create

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

func CreateRedisConnection(ctx context.Context) (*redis.Client, error){
	addr := os.Getenv("REDIS_ADDR")
	if addr == ""{
		addr = "6379"
	}

	pwd := os.Getenv("REDIS_PASSWORD")
	db:= 0

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: pwd,
		DB: db,

		DialTimeout:  10 * time.Second,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        PoolSize:     10, 
        MinIdleConns: 2,
	})	

	if err := rdb.Ping(ctx).Err(); err != nil{
		rdb.Close()
		return nil, fmt.Errorf("Failed to connect to Redis: %w", err)
	}
	return rdb, nil
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
	names := []string{"menu", "version", "item", "users", "restaurant", "session", "order_items"}
	schemas := []string{
		"id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, restaurant_id INT NOT NULL REFERENCES restaurant(id), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()",
		"id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, menu_id INT NOT NULL REFERENCES menu(id), s3_url TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), is_active BOOLEAN NOT NULL",
		"id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, version_id INT NOT NULL REFERENCES version(id), name TEXT NOT NULL, description TEXT, price NUMERIC(10,2) NOT NULL, category TEXT, modifiers JSONB",
		"id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, uuid TEXT NOT NULL UNIQUE, email TEXT, name TEXT, password TEXT",
		"id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, name TEXT NOT NULL, address TEXT, phone TEXT",
		"id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, host INT NOT NULL REFERENCES users(id), status TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()",
		"id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, session_id INT NOT NULL REFERENCES session(id), user_id INT NOT NULL REFERENCES users(id), item TEXT NOT NULL, modifiers JSONB, price INT NOT NULL",
	}

	for i, name := range names {
		schema := schemas[i]
		fmt.Printf("Creating table %s with schema %s\n", name, schemas)
		rows, err := CreateTable(conn, name, schema)
		if err != nil {
			log.Printf("Error creating table %s: %v\n", name, err)
			continue
		}
		fmt.Printf("Created %d rows\n", rows)
	}

	return nil
}
