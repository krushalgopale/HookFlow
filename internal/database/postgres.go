package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		return nil, fmt.Errorf("POSTGRES_URL is not set in .env")
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
