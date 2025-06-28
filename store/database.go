package store

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
)

func Open() (*sql.DB, error) {
	connectionStr := os.Getenv("DATABASE_URL")

	if connectionStr == "" {
		return nil, errors.New("DATABASE_URL environment variable is not set")
	}

	DB, err := sql.Open("postgres", connectionStr)

	if err != nil {
		return nil, err
	}

	if err := DB.Ping(); err != nil {
		return nil, err
	}

	return DB, nil
}
