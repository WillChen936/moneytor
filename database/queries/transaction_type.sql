-- name: GetTransactionType :one
SELECT * 
  FROM transaction_types
 WHERE id = $1 LIMIT 1;

-- name: ListTransactionTypes :many
SELECT * 
  FROM transaction_types
 ORDER BY id;