package database

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {

	databaseURL, exists := os.LookupEnv("Database_URL")

	if !exists || databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
