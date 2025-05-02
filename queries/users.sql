-- name: CreateUser :one
INSERT INTO users (name, email, password_hashed)
VALUES ($1, $2, $3) RETURNING id;

-- name: CreateSession :one
INSERT INTO tokens (user_id, issued_at, expires_at)
VALUES($1, $2, $3) RETURNING id;

-- name: FindUserByNameOrEmail :one
SELECT id, password_hashed FROM users
WHERE name=sqlc.arg(name_or_email) OR email=sqlc.arg(name_or_email);

-- name: GetUserById :one
SELECT name, created_at, last_login FROM users WHERE id=$1;

-- name: DeleteById :one
DELETE FROM tokens WHERE id=$1 RETURNING user_id, expires_at;