-- name: CreateEntry :one
INSERT INTO entries (
     account_id  
    ,category_id 
    ,amount      
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListEntries :many
SELECT *
  FROM entries
 WHERE account_id = $1
 ORDER BY id DESC
 LIMIT $2
 OFFSET $3;