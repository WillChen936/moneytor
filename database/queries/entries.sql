-- name: CreateEntry :one
INSERT INTO entries (
     name
    ,note
    ,account_id  
    ,category_id 
    ,amount      
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetEntry :one
SELECT *
  FROM entries
 WHERE id = $1
 LIMIT 1;

-- name: ListEntries :many
SELECT *
  FROM entries
 ORDER BY id DESC
 LIMIT $1
 OFFSET $2;

-- name: ListEntriesByAccountID :many
SELECT *
  FROM entries
 WHERE account_id = $1
 ORDER BY id DESC
 LIMIT $2
 OFFSET $3;