-- name: CreateUser :one
INSERT INTO users (
  email, 
  hashed_pass,
  username, 
  full_name
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET hashed_pass = COALESCE(sqlc.narg(hashed_pass), hashed_pass),
    email = COALESCE(sqlc.narg(email), email),
    full_name = COALESCE(sqlc.narg(full_name), full_name)
WHERE username = @username
RETURNING *;