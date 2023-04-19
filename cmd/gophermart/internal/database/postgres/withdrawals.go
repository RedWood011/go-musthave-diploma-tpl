package postgres

import (
	"context"
	"time"

	"github.com/RedWood011/cmd/gophermart/internal/apperrors"
	"github.com/RedWood011/cmd/gophermart/internal/entity"
	"github.com/jackc/pgx/v4"
)

func (r *Repository) SaveWithdraw(ctx context.Context, withdraw entity.Withdraw, userID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var enough bool
	row := tx.QueryRow(ctx, "SELECT (balance >= $1) FROM users WHERE id = $2 FOR UPDATE", withdraw.Sum, userID)
	err = row.Scan(&enough)
	if err != nil {
		return err
	}
	if !enough {
		return apperrors.ErrNoMoney
	}
	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $1, spend = spend + $1 WHERE id = $2", withdraw.Sum, userID)
	if err != nil {
		return err
	}
	sqlCreateWithdraw := `INSERT INTO withdrawals (order_number, user_id, sum, processed_at) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ctx, sqlCreateWithdraw, withdraw.OrderNumber, userID, withdraw.Sum, time.Now())
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) GetWithdrawals(ctx context.Context, userID string) ([]entity.Withdraw, error) {
	var result []entity.Withdraw
	sqlGetOrders := `SELECT order_number, sum, processed_at FROM withdrawals
					 WHERE user_id = $1 ORDER BY processed_at ASC`
	rows, err := r.db.Query(ctx, sqlGetOrders, userID)
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var res entity.Withdraw
		if err = rows.Scan(&res.UserID, &res.OrderNumber, &res.Sum, &res.ProcessedAt); err != nil {
			return result, nil
		}
		result = append(result, res)

	}
	return result, nil
}
