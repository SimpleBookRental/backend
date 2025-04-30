-- name: GetRental :one
SELECT r.*, 
       u.username as user_username, 
       b.title as book_title, 
       b.author as book_author
FROM rentals r
JOIN users u ON r.user_id = u.id
JOIN books b ON r.book_id = b.id
WHERE r.id = $1 LIMIT 1;

-- name: ListRentals :many
SELECT r.*, 
       u.username as user_username, 
       b.title as book_title, 
       b.author as book_author
FROM rentals r
JOIN users u ON r.user_id = u.id
JOIN books b ON r.book_id = b.id
ORDER BY r.rental_date DESC
LIMIT $1
OFFSET $2;

-- name: ListRentalsByUser :many
SELECT r.*, 
       u.username as user_username, 
       b.title as book_title, 
       b.author as book_author
FROM rentals r
JOIN users u ON r.user_id = u.id
JOIN books b ON r.book_id = b.id
WHERE r.user_id = $1
ORDER BY r.rental_date DESC
LIMIT $2
OFFSET $3;

-- name: ListRentalsByBook :many
SELECT r.*, 
       u.username as user_username, 
       b.title as book_title, 
       b.author as book_author
FROM rentals r
JOIN users u ON r.user_id = u.id
JOIN books b ON r.book_id = b.id
WHERE r.book_id = $1
ORDER BY r.rental_date DESC
LIMIT $2
OFFSET $3;

-- name: ListActiveRentals :many
SELECT r.*, 
       u.username as user_username, 
       b.title as book_title, 
       b.author as book_author
FROM rentals r
JOIN users u ON r.user_id = u.id
JOIN books b ON r.book_id = b.id
WHERE r.status = 'active'
ORDER BY r.due_date ASC
LIMIT $1
OFFSET $2;

-- name: ListOverdueRentals :many
SELECT r.*, 
       u.username as user_username, 
       b.title as book_title, 
       b.author as book_author
FROM rentals r
JOIN users u ON r.user_id = u.id
JOIN books b ON r.book_id = b.id
WHERE r.status = 'active' AND r.due_date < NOW()
ORDER BY r.due_date ASC
LIMIT $1
OFFSET $2;

-- name: CreateRental :one
INSERT INTO rentals (
  user_id,
  book_id,
  rental_date,
  due_date,
  status
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateRentalStatus :one
UPDATE rentals
SET 
  status = $2,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ReturnRental :one
UPDATE rentals
SET 
  return_date = NOW(),
  status = 'returned',
  updated_at = NOW()
WHERE id = $1 AND status = 'active'
RETURNING *;

-- name: ExtendRental :one
UPDATE rentals
SET 
  due_date = $2,
  updated_at = NOW()
WHERE id = $1 AND status = 'active'
RETURNING *;

-- name: DeleteRental :exec
DELETE FROM rentals
WHERE id = $1;
