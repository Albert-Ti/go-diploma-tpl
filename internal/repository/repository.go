package repository

import (
	"context"

	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository interface {
	RegisterTx(ctx context.Context, login string, pass string) (int, error)
	GetUser(ctx context.Context, login string) (int, string, error)
	CreateOrder(ctx context.Context, order string, userID int) (int, error)
	GetOrders(ctx context.Context, userID int) ([]models.OrdersResp, error)
	GetBalance(ctx context.Context, userID int) (models.BalanceResp, error)
	UpdateStatusOrder(ctx context.Context, order string, status models.OrderStatus) error
	ProcessedOrderTx(ctx context.Context, order string, status models.OrderStatus, accrual float64, userID int) error
	BalanceWithdrawalsTx(ctx context.Context, order string, sum float64, userID int) error
	GetWithdrawals(ctx context.Context, userID int) ([]models.WithdrawalsResp, error)
}

func NewRepository(dsn string) (Repository, error) {
	return NewPgStorage(dsn)
}
