-- name: CreateCategory :one
INSERT INTO categories (
    name, 
    transaction_type_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetCategory :one
SELECT *
  FROM categories
 WHERE id = $1
 LIMIT 1;

-- name: ListCategories :many
SELECT * 
  FROM categories
 ORDER BY id
 LIMIT $1
OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
   SET name = COALESCE(sqlc.narg('name'), name),
       transaction_type_id = COALESCE(sqlc.narg('transaction_type_id'), transaction_type_id),
       updated_at = NOW()
 WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
 WHERE id = $1;