-- name: CreatePayment :one
INSERT INTO payments (
  from_trader_id,
  to_trader_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payments
WHERE 
    from_trader_id = $1 OR
    to_trader_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;
