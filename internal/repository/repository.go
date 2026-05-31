package repository

import (
	"context"
	"log/slog"

	"github.com/Albert-Ti/go-diploma-tpl/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository interface {
	AddUser(ctx context.Context, login string, pass string) error
	GetUser(ctx context.Context, login string) (int, string, error)
}

func NewRepository() (Repository, error) {
	slog.Info("Using database storage")

	m, err := migrate.New("file://migrations", config.Envs.DatabaseURI)
	if err != nil {
		return nil, err
	}
	defer m.Close()

	err = m.Up()
	switch err {
	case nil:
		slog.Info("Migrations have been successfully applied")
	case migrate.ErrNoChange:
		slog.Info("The database schema is up-to-date and no migrations are required")
	default:
		slog.Error("Migrations failed", "error", err)
		return nil, err
	}

	return NewPgStorage(config.Envs.DatabaseURI)
}
