package postgres_storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/redblood-pixel/pastebin/internal/domain"
)

const createUserQuery = `INSERT INTO users (name, email, password_hashed)
VALUES ($1, $2, $3) RETURNING id;`

const createSessionQuery = `INSERT INTO tokens (user_id, issued_at, expires_at)
VALUES($1, $2, $3) RETURNING id;`

const findUserByNameOrEmailQuery = `SELECT id, password_hashed FROM users
WHERE name=$1 OR email=$1;`

const getUserByIdQuery = `SELECT name, created_at, last_login FROM users WHERE id=$1;`

const deleteByIdQuery = `DELETE FROM tokens WHERE id=$1 RETURNING user_id, expires_at;`

func (r *PostgresStorage) CreateUser(
	ctx context.Context,
	tx pgx.Tx,
	name, email, password_hashed string,
) (int, error) {
	var id int
	row := tx.QueryRow(ctx, createUserQuery, name, email, password_hashed)
	err := row.Scan(&id)

	return id, err
}

func (r *PostgresStorage) CreateSession(
	ctx context.Context,
	tx pgx.Tx,
	userID int, expiryTime time.Duration,
) (uuid.UUID, error) {
	var id uuid.UUID
	row := tx.QueryRow(ctx, createSessionQuery, userID, time.Now(), time.Now().Add(expiryTime))
	if err := row.Scan(&id); err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (r *PostgresStorage) FindUserByNameOrEmail(
	ctx context.Context,
	nameOrEmail string,
) (int, string, error) {
	var (
		userID         int
		passwordHashed string
	)
	row := r.db.Pool.QueryRow(ctx, findUserByNameOrEmailQuery, nameOrEmail)
	if err := row.Scan(&userID, &passwordHashed); err != nil {
		return 0, "", err
	}
	fmt.Println(userID, passwordHashed)
	return userID, passwordHashed, nil
}

func (r *PostgresStorage) GetUserById(
	ctx context.Context,
	userID int,
) (domain.User, error) {
	var user domain.User
	row := r.db.Pool.QueryRow(ctx, getUserByIdQuery, userID)
	if err := row.Scan(&user.Name, &user.CreatedAt, &user.LastLogin); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, domain.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}

func (r *PostgresStorage) DeleteSessionById(
	ctx context.Context,
	tx pgx.Tx,
	sessionID uuid.UUID,
) (int, time.Time, error) {
	var (
		userID   int
		expireAt time.Time
	)
	row := tx.QueryRow(ctx, deleteByIdQuery, sessionID)
	if err := row.Scan(&userID, &expireAt); err != nil {
		return 0, time.Now(), err
	}
	fmt.Println(userID, userID)
	return userID, expireAt, nil
}
