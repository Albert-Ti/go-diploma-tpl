package repository

import (
	"context"
	"errors"

	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrConflictuser error = errors.New("Order has already been uploaded by this user")

type PgStorage struct {
	pool          *pgxpool.Pool
	pgErrConflict error
}

func NewPgStorage(dsn string) (*PgStorage, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &PgStorage{
		pool:          pool,
		pgErrConflict: errors.New("Order is already exist"),
	}, nil
}

func (pg *PgStorage) addUserWithTx(ctx context.Context, tx pgx.Tx, login string, pass string) (int, error) {
	sql := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`
	var id int
	err := tx.QueryRow(ctx, sql, login, pass).Scan(&id)
	return id, err
}

func (pg *PgStorage) addUserBalanceWithTx(ctx context.Context, tx pgx.Tx, userID int) error {
	sql := `INSERT INTO user_balance (user_id) VALUES ($1)`
	_, err := tx.Exec(ctx, sql, userID)
	return err
}

func (pg *PgStorage) Registration(ctx context.Context, login string, pass string) (int, error) {
	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	userID, err := pg.addUserWithTx(ctx, tx, login, pass)
	if err != nil {
		return 0, err
	}

	if err := pg.addUserBalanceWithTx(ctx, tx, userID); err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return userID, nil
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

func (pg *PgStorage) CreateOrder(ctx context.Context, order string, userID int) error {
	sql := `
		INSERT INTO orders (number, status, user_id)
		VALUES ($1, $2, $3)
	`

	_, err := pg.pool.Exec(ctx, sql, order, StatusNew, userID)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PgStorage) GetOrders(ctx context.Context, userID int) ([]models.GetOrdersResp, error) {
	sql := `
		SELECT number, status, accrual, uploaded_at FROM orders
		WHERE user_id = $1
	`

	rows, err := pg.pool.Query(ctx, sql, userID)
	if err != nil {
		return []models.GetOrdersResp{}, err
	}

	results := []models.GetOrdersResp{}

	for rows.Next() {
		var m models.GetOrdersResp
		err := rows.Scan(&m.Number, &m.Status, &m.Accrual, &m.UploadedAt)
		if err != nil {
			return []models.GetOrdersResp{}, err
		}

		results = append(results, m)
	}

	return results, nil
}

func (pg *PgStorage) GetOrderUserID(ctx context.Context, order string) (int, error) {
	sql := `SELECT user_id FROM orders WHERE number = $1`
	var userID int
	err := pg.pool.QueryRow(ctx, sql, order).Scan(&userID)
	return userID, err
}
