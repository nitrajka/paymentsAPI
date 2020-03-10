-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1;

-- name: GetPayments :many
SELECT * FROM payments;

-- name: CreatePayment :one
INSERT INTO payments (amount, description, sender, datetime)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetBalance :one
SELECT sum(amount) FROM payments;