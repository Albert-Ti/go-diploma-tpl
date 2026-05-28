package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgStorage struct {
	pool *pgxpool.Pool
}

func NewPgStorage(dsn string) (*PgStorage, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &PgStorage{
		pool: pool,
	}, nil
}
