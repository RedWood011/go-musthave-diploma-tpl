package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/RedWood011/cmd/gophermart/internal/entity"
	"github.com/jackc/pgx/v4"
)

type order struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Number     string    `json:"number"`
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual"`
}

func (o order) toEntity() entity.Order {
	return entity.Order{
		ID:         o.ID,
		UserID:     o.UserID,
		Number:     o.Number,
		UploadedAt: o.UploadedAt,
		Status:     o.Status,
		Accrual:    o.Accrual,
	}
}

func (r *Repository) SaveOrder(ctx context.Context, order entity.Order) error {
	sqlCreateOrder := `INSERT INTO orders (user_id, number, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4,$5)`
	_, err := r.db.Exec(ctx, sqlCreateOrder, order.UserID, order.Number, order.Status, order.Accrual, order.UploadedAt)
	return err
}

func (r *Repository) GetOrder(ctx context.Context, orderNum string) (entity.Order, error) {
	var res order
	queryGetOrder := `SELECT id, user_id, number, status, accrual, uploaded_at FROM orders WHERE number = $1`
	result := r.db.QueryRow(ctx, queryGetOrder, orderNum)
	if err := result.Scan(&res.ID, &res.UserID, &res.Number, &res.Status, &res.Accrual, &res.UploadedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, err
		}

		return entity.Order{}, err
	}

	return res.toEntity(), nil

}

func (r *Repository) GetUserOrders(ctx context.Context, userUID string) ([]entity.Order, error) {
	var result []entity.Order
	queryGetOrders := `SELECT number, status, uploaded_at, accrual FROM orders
					 WHERE user_id = $1 ORDER BY uploaded_at`
	rows, err := r.db.Query(ctx, queryGetOrders, userUID)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var res order
		err = rows.Scan(&res.Number, &res.Status, &res.UploadedAt, &res.Accrual)
		if err != nil {
			return nil, err
		}
		result = append(result, res.toEntity())
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, err
}

func (r *Repository) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	var result []entity.Order
	query := " select user_id, number, uploaded_at, status, accrual  from orders where status IN ('NEW','PROCESSING') "
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var res order
		err = rows.Scan(&res.UserID, &res.Number, &res.UploadedAt, &res.Status, &res.Accrual)
		if err != nil {
			return nil, err
		}
		result = append(result, res.toEntity())
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) UpdateOrders(ctx context.Context, orders []entity.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	query := `UPDATE orders SET accrual = $1, status = $2 WHERE number = $3`
	defer tx.Rollback(ctx)
	for _, value := range orders {
		_, err = tx.Exec(ctx, query, value.Accrual, value.Status, value.Number)
		if err != nil {
			return err
		}
		tx.Commit(ctx)
	}
	return nil
}
