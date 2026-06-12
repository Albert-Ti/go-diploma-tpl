package repository

import (
	"context"
	"log/slog"

	"github.com/Albert-Ti/go-diploma-tpl/internal/config"
	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository interface {
	Register(ctx context.Context, login string, pass string) (int, error)
	GetUser(ctx context.Context, login string) (int, string, error)
	CreateOrder(ctx context.Context, order string, userID int) error
	GetOrders(ctx context.Context, userID int) ([]models.GetOrdersResp, error)
	GetOrderUserID(ctx context.Context, order string) (int, error)
	GetBalance(ctx context.Context, userID int) (models.GetBalanceResp, error)
	UpdateStatusOrder(ctx context.Context, order string, status models.OrderStatus) error
	ProcessedOrder(ctx context.Context, order string, status models.OrderStatus, accrual float64, userID int) error
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
