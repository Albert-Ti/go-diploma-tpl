package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

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

func (pg *PgStorage) RegisterTx(ctx context.Context, login string, pass string) (int, error) {
	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback(ctx)

	var userID int
	errUser := tx.QueryRow(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`,
		login, pass).Scan(&userID)
	if errUser != nil {
		return 0, errUser
	}

	_, errBalance := tx.Exec(ctx,
		`INSERT INTO user_balance (user_id) VALUES ($1)`,
		userID)
	if errBalance != nil {
		return 0, errBalance
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return userID, nil
}

func (pg *PgStorage) GetUser(ctx context.Context, login string) (int, string, error) {
	sql := `SELECT id, password FROM users WHERE login = $1`

	row := pg.pool.QueryRow(ctx, sql, login)

	var userID int
	var password string

	err := row.Scan(&userID, &password)
	if err != nil {
		return 0, "", err
	}
	return userID, password, nil
}

func (pg *PgStorage) CreateOrder(ctx context.Context, order string, userID int) (int, error) {
	sql := `
        INSERT INTO orders (number, status, user_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (number) DO UPDATE 
        SET number = EXCLUDED.number
        RETURNING user_id, (xmax = 0) AS is_inserted;
    `

	var existingUserID int
	var isNew bool

	err := pg.pool.QueryRow(ctx, sql, order, models.StatusNew, userID).
		Scan(&existingUserID, &isNew)
	if err != nil {
		return 0, err
	}

	if isNew {
		return 0, nil
	}

	return existingUserID, err
}

func (pg *PgStorage) GetOrders(ctx context.Context, userID int) ([]models.OrdersResp, error) {
	sql := `
		SELECT number, status, accrual, uploaded_at FROM orders
		WHERE user_id = $1
	`
	rows, err := pg.pool.Query(ctx, sql, userID)
	if err != nil {
		return []models.OrdersResp{}, err
	}

	results := []models.OrdersResp{}

	for rows.Next() {
		var m models.OrdersResp
		err := rows.Scan(&m.Number, &m.Status, &m.Accrual, &m.UploadedAt)
		if err != nil {
			return []models.OrdersResp{}, err
		}

		results = append(results, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (pg *PgStorage) UpdateStatusOrder(ctx context.Context, order string, status models.OrderStatus) error {
	sql := `UPDATE orders SET status = $1 WHERE number = $2`

	_, err := pg.pool.Exec(ctx, sql, status, order)

	return err
}

func (pg *PgStorage) ProcessedOrderTx(ctx context.Context, order string, status models.OrderStatus, accrual float64, userID int) error {
	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`,
		status, accrual, order)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE user_balance SET current = current + $1 WHERE user_id = $2`,
		accrual, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (pg *PgStorage) GetBalance(ctx context.Context, userID int) (models.BalanceResp, error) {
	sql := `SELECT current, withdrawn FROM user_balance WHERE user_id = $1`

	var m models.BalanceResp
	if err := pg.pool.QueryRow(ctx, sql, userID).Scan(&m.Current, &m.Withdraw); err != nil {
		return models.BalanceResp{}, err
	}

	return m, nil
}

func (pg *PgStorage) BalanceWithdrawalsTx(ctx context.Context, order string, sum float64, userID int) error {
	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	//FOR UPDATE — это инструкция, которая блокирует выбранные строки,
	// чтобы другие транзакции не могли их изменить, пока текущая транзакция не завершится.
	var currentBalance float64
	err = tx.QueryRow(ctx,
		`SELECT current FROM user_balance WHERE user_id = $1 FOR UPDATE`,
		userID).Scan(&currentBalance)
	if err != nil {
		return err
	}

	if currentBalance < sum {
		return fmt.Errorf("%w: current balance %.2f, requested %.2f", ErrInsufficientFunds, currentBalance, sum)
	}

	_, errBalance := tx.Exec(ctx,
		`UPDATE user_balance SET current = current - $1, withdrawn = withdrawn + $1 WHERE user_id = $2`,
		sum, userID)
	if errBalance != nil {
		return errBalance
	}

	_, errWithdrawals := tx.Exec(ctx,
		`INSERT INTO withdrawals (order_number, sum, user_id) VALUES($1, $2, $3)`,
		order, sum, userID)
	if errWithdrawals != nil {
		return errWithdrawals
	}

	return tx.Commit(ctx)
}

func (pg *PgStorage) GetWithdrawals(ctx context.Context, userID int) ([]models.WithdrawalsResp, error) {
	sql := `SELECT order_number, sum, processed_at FROM withdrawals WHERE user_id = $1`

	rows, err := pg.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}

	results := make([]models.WithdrawalsResp, 0)
	for rows.Next() {
		var m models.WithdrawalsResp
		err := rows.Scan(&m.Order, &m.Sum, &m.ProcessedAt)
		if err != nil {
			return nil, err
		}

		results = append(results, m)
	}

	return results, nil
}
