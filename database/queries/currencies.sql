-- name: GetCurrency :one
SELECT * 
  FROM currencies
 WHERE id = $1 LIMIT 1;

-- name: ListCurrencies :many
SELECT * 
  FROM currencies
 ORDER BY id;