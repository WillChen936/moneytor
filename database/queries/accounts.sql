-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    currency_id,
    balance    
) VALUES (
  $1, $2, sqlc.narg('balance')
)
RETURNING *;

-- name: GetAccount :one
SELECT *
  FROM accounts
 WHERE id = $1
 LIMIT 1;

-- name: UpdateAccount :one
UPDATE accounts
   SET owner = COALESCE(sqlc.narg('owner'), owner),
       balance = COALESCE(sqlc.narg('balance'), balance),
       updated_at = NOW()
 WHERE id = $1
RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE accounts
   SET balance = balance + sqlc.narg('amount'),
       updated_at = NOW()
 WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
 WHERE id = $1;