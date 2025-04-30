-- name: GetPayment :one
SELECT p.*, 
       u.username as user_username,
       r.id as rental_id,
       b.title as book_title
FROM payments p
JOIN users u ON p.user_id = u.id
LEFT JOIN rentals r ON p.rental_id = r.id
LEFT JOIN books b ON r.book_id = b.id
WHERE p.id = $1 LIMIT 1;

-- name: ListPayments :many
SELECT p.*, 
       u.username as user_username,
       r.id as rental_id,
       b.title as book_title
FROM payments p
JOIN users u ON p.user_id = u.id
LEFT JOIN rentals r ON p.rental_id = r.id
LEFT JOIN books b ON r.book_id = b.id
ORDER BY p.payment_date DESC
LIMIT $1
OFFSET $2;

-- name: ListPaymentsByUser :many
SELECT p.*, 
       u.username as user_username,
       r.id as rental_id,
       b.title as book_title
FROM payments p
JOIN users u ON p.user_id = u.id
LEFT JOIN rentals r ON p.rental_id = r.id
LEFT JOIN books b ON r.book_id = b.id
WHERE p.user_id = $1
ORDER BY p.payment_date DESC
LIMIT $2
OFFSET $3;

-- name: ListPaymentsByRental :many
SELECT p.*, 
       u.username as user_username,
       r.id as rental_id,
       b.title as book_title
FROM payments p
JOIN users u ON p.user_id = u.id
LEFT JOIN rentals r ON p.rental_id = r.id
LEFT JOIN books b ON r.book_id = b.id
WHERE p.rental_id = $1
ORDER BY p.payment_date DESC;

-- name: CreatePayment :one
INSERT INTO payments (
  user_id,
  rental_id,
  amount,
  payment_date,
  payment_method,
  status,
  transaction_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET 
  status = $2,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;

-- name: GetRevenueReport :many
SELECT 
  DATE_TRUNC('month', payment_date) as month,
  SUM(amount) as total_revenue,
  COUNT(*) as payment_count
FROM payments
WHERE status = 'completed'
  AND payment_date BETWEEN $1 AND $2
GROUP BY DATE_TRUNC('month', payment_date)
ORDER BY month;
