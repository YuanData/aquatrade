// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: trader.sql

package db

import (
	"context"
)

const addTraderBalance = `-- name: AddTraderBalance :one
UPDATE traders
SET balance = balance + $1
WHERE id = $2
RETURNING id, holder, balance, currency, created_at
`

type AddTraderBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddTraderBalance(ctx context.Context, arg AddTraderBalanceParams) (Trader, error) {
	row := q.db.QueryRowContext(ctx, addTraderBalance, arg.Amount, arg.ID)
	var i Trader
	err := row.Scan(
		&i.ID,
		&i.Holder,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const createTrader = `-- name: CreateTrader :one
INSERT INTO traders (
  holder,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING id, holder, balance, currency, created_at
`

type CreateTraderParams struct {
	Holder   string `json:"holder"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateTrader(ctx context.Context, arg CreateTraderParams) (Trader, error) {
	row := q.db.QueryRowContext(ctx, createTrader, arg.Holder, arg.Balance, arg.Currency)
	var i Trader
	err := row.Scan(
		&i.ID,
		&i.Holder,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTrader = `-- name: DeleteTrader :exec
DELETE FROM traders
WHERE id = $1
`

func (q *Queries) DeleteTrader(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTrader, id)
	return err
}

const getTrader = `-- name: GetTrader :one
SELECT id, holder, balance, currency, created_at FROM traders
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTrader(ctx context.Context, id int64) (Trader, error) {
	row := q.db.QueryRowContext(ctx, getTrader, id)
	var i Trader
	err := row.Scan(
		&i.ID,
		&i.Holder,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const listTraders = `-- name: ListTraders :many
SELECT id, holder, balance, currency, created_at FROM traders
WHERE holder = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListTradersParams struct {
	Holder string `json:"holder"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListTraders(ctx context.Context, arg ListTradersParams) ([]Trader, error) {
	rows, err := q.db.QueryContext(ctx, listTraders, arg.Holder, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Trader
	for rows.Next() {
		var i Trader
		if err := rows.Scan(
			&i.ID,
			&i.Holder,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTrader = `-- name: UpdateTrader :one
UPDATE traders
SET balance = $2
WHERE id = $1
RETURNING id, holder, balance, currency, created_at
`

type UpdateTraderParams struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func (q *Queries) UpdateTrader(ctx context.Context, arg UpdateTraderParams) (Trader, error) {
	row := q.db.QueryRowContext(ctx, updateTrader, arg.ID, arg.Balance)
	var i Trader
	err := row.Scan(
		&i.ID,
		&i.Holder,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}
