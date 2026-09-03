package database

import (
	"context"
	"log"
	"time"

	"gin-quickstart/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {

	cfg := config.Load()

	config, err := pgxpool.ParseConfig(cfg.Database_URL)
	if err != nil {
		log.Printf("DB URL config error%v", err)
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

	log.Println("INFO: DB started successfully")
	return pool, nil
}
