-- name: GetCategory :one
SELECT *
  FROM categories
 WHERE id = $1
 LIMIT 1;

-- name: CreateCategory :one
INSERT INTO categories (
    name, 
    transaction_type_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateCategory :exec
UPDATE categories
   SET name = $2,
       transaction_type_id = $3
 WHERE id = $1;

-- name: DeleteCategory :exec
DELETE FROM categories
 WHERE id = $1;