-- name: CreateUser :one
INSERT INTO users (name, email, password_hashed) VALUES ($1, $2, $3) RETURNING id;