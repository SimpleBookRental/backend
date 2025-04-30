-- name: GetBook :one
SELECT * FROM books
WHERE id = $1 LIMIT 1;

-- name: GetBookByISBN :one
SELECT * FROM books
WHERE isbn = $1 LIMIT 1;

-- name: ListBooks :many
SELECT b.*, c.name as category_name
FROM books b
LEFT JOIN categories c ON b.category_id = c.id
ORDER BY b.title
LIMIT $1
OFFSET $2;

-- name: ListBooksByCategory :many
SELECT b.*, c.name as category_name
FROM books b
JOIN categories c ON b.category_id = c.id
WHERE b.category_id = $1
ORDER BY b.title
LIMIT $2
OFFSET $3;

-- name: SearchBooks :many
SELECT b.*, c.name as category_name
FROM books b
LEFT JOIN categories c ON b.category_id = c.id
WHERE 
  ($1::text IS NULL OR b.title ILIKE '%' || $1 || '%') AND
  ($2::text IS NULL OR b.author ILIKE '%' || $2 || '%') AND
  ($3::text IS NULL OR b.isbn = $3) AND
  ($4::int IS NULL OR b.published_year = $4) AND
  ($5::int IS NULL OR b.category_id = $5) AND
  ($6::boolean IS NULL OR ($6 = true AND b.available_copies > 0) OR ($6 = false))
ORDER BY b.title
LIMIT $7
OFFSET $8;

-- name: CreateBook :one
INSERT INTO books (
  title,
  author,
  isbn,
  description,
  published_year,
  publisher,
  total_copies,
  available_copies,
  category_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $7, $8
)
RETURNING *;

-- name: UpdateBook :one
UPDATE books
SET 
  title = COALESCE(sqlc.narg(title), title),
  author = COALESCE(sqlc.narg(author), author),
  isbn = COALESCE(sqlc.narg(isbn), isbn),
  description = COALESCE(sqlc.narg(description), description),
  published_year = COALESCE(sqlc.narg(published_year), published_year),
  publisher = COALESCE(sqlc.narg(publisher), publisher),
  category_id = COALESCE(sqlc.narg(category_id), category_id),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateBookCopies :one
UPDATE books
SET 
  total_copies = COALESCE(sqlc.narg(total_copies), total_copies),
  available_copies = COALESCE(sqlc.narg(available_copies), available_copies),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DecrementAvailableCopies :one
UPDATE books
SET 
  available_copies = available_copies - 1,
  updated_at = NOW()
WHERE id = $1 AND available_copies > 0
RETURNING *;

-- name: IncrementAvailableCopies :one
UPDATE books
SET 
  available_copies = available_copies + 1,
  updated_at = NOW()
WHERE id = $1 AND available_copies < total_copies
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;
