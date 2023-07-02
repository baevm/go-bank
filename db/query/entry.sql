-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 AND account_id = $2 
LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $3
ORDER BY id
LIMIT $1
OFFSET $2;
