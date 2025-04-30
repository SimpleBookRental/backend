-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoryByName :one
SELECT * FROM categories
WHERE name = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: ListAllCategories :many
SELECT * FROM categories
ORDER BY name;

-- name: CreateCategory :one
INSERT INTO categories (
  name,
  description
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET 
  name = COALESCE(sqlc.narg(name), name),
  description = COALESCE(sqlc.narg(description), description),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;
