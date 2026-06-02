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

func (pg *PgStorage) AddUser(ctx context.Context, login string, pass string) (int, error) {
	sql := `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id
	`
	var id int
	err := pg.pool.QueryRow(ctx, sql, login, pass).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (pg *PgStorage) GetUser(ctx context.Context, login string) (int, string, error) {
	sql := `
		SELECT id, password FROM users WHERE login = $1
	`
	row := pg.pool.QueryRow(ctx, sql, login)

	var userID int
	var password string

	err := row.Scan(&userID, &password)
	if err != nil {
		return 0, "", err
	}
	return userID, password, nil
}
