-- name: CreateAccount :one
INSERT INTO accounts (
    user_id,
    name,
    currency_id,
    balance
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetAccount :one
SELECT *
  FROM accounts
 WHERE id = $1
   AND user_id = $2
 LIMIT 1;

-- name: ListAccounts :many
SELECT *
  FROM accounts
 WHERE user_id = $1
 ORDER BY id
 LIMIT $2
OFFSET $3;

-- name: UpdateAccountBalance :one
UPDATE accounts
   SET balance = balance + sqlc.arg('amount'),
       updated_at = NOW()
 WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
 WHERE id = $1
   AND user_id = $2;
