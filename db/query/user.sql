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
    password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
    email = COALESCE(sqlc.narg(email), email),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    is_email_activated = COALESCE(sqlc.narg(is_email_activated), is_email_activated)
WHERE username = @username
RETURNING *;