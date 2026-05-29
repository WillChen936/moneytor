-- name: CreateCategory :one
INSERT INTO categories (
    user_id,
    name,
    transaction_type_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetCategory :one
SELECT *
  FROM categories
 WHERE id = $1
   AND user_id = $2
 LIMIT 1;

-- name: ListCategories :many
SELECT *
  FROM categories
 WHERE user_id = $1
 ORDER BY id
 LIMIT $2
OFFSET $3;

-- name: UpdateCategory :one
UPDATE categories
   SET name = COALESCE(sqlc.narg('name'), name),
       transaction_type_id = COALESCE(sqlc.narg('transaction_type_id'), transaction_type_id),
       updated_at = NOW()
 WHERE id = $1
   AND user_id = sqlc.arg('user_id')
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
 WHERE id = $1
   AND user_id = $2;
