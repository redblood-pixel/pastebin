-- name: CreateUser :one
INSERT INTO users (name, email, password_hashed) VALUES ($1, $2, $3) RETURNING id;

-- name: CreateSession :one
INSERT INTO tokens (user_id, issued_at, expires_at) VALUES($1, $2, $3) RETURNING id;