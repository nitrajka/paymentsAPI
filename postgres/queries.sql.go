// Code generated by sqlc. DO NOT EDIT.
// source: queries.sql

package postgres

import (
	"context"
	"time"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payments (amount, description, sender, datetime)
VALUES ($1, $2, $3, $4) RETURNING id, amount, description, sender, datetime
`

type CreatePaymentParams struct {
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Sender      string    `json:"sender"`
	Datetime    time.Time `json:"datetime"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRowContext(ctx, createPayment,
		arg.Amount,
		arg.Description,
		arg.Sender,
		arg.Datetime,
	)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Description,
		&i.Sender,
		&i.Datetime,
	)
	return i, err
}

const getBalance = `-- name: GetBalance :one
SELECT sum(amount) FROM payments
`

func (q *Queries) GetBalance(ctx context.Context) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getBalance)
	var sum interface{}
	err := row.Scan(&sum)
	return sum, err
}

const getPayment = `-- name: GetPayment :one
SELECT id, amount, description, sender, datetime FROM payments
WHERE id = $1
`

func (q *Queries) GetPayment(ctx context.Context, id int32) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getPayment, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Description,
		&i.Sender,
		&i.Datetime,
	)
	return i, err
}

const getPayments = `-- name: GetPayments :many
SELECT id, amount, description, sender, datetime FROM payments
`

func (q *Queries) GetPayments(ctx context.Context) ([]Payment, error) {
	rows, err := q.db.QueryContext(ctx, getPayments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Payment
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.ID,
			&i.Amount,
			&i.Description,
			&i.Sender,
			&i.Datetime,
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
