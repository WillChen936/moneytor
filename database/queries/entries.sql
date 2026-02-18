-- name: CreateEntry :one
INSERT INTO entries (
     name
    ,note
    ,from_account_id  
    ,to_account_id  
    ,category_id 
    ,amount      
) VALUES (
  $1, $2, $3, $4, $5, $6
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
 WHERE from_account_id = sqlc.arg('account_id') OR to_account_id = sqlc.arg('account_id')
 ORDER BY id DESC
 LIMIT $1
 OFFSET $2;