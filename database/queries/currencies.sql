-- name: GetCurrency :one
SELECT * 
  FROM currencies
 WHERE id = $1 LIMIT 1;