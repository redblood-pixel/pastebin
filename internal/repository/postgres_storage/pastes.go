package postgres_storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/redblood-pixel/pastebin/internal/domain"
)

const createPasteQuery = `INSERT INTO pastes(
	title,
	expires_at,
	visibility,
	user_id
) VALUES ($1, $2, $3, $4) RETURNING id;`

const getPasteByIDQuery = `SELECT
	title,
	created_at,
	expires_at,
	visibility,
	last_visited,
	user_id
FROM pastes WHERE id=$1;
`

const getUsersPastesQuery = `SELECT
	id,
	title,
	created_at,
	expires_at,
	visibility,
	last_visited
FROM pastes WHERE user_id=$1;
`

const updateLastVisitedQuery = `UPDATE pastes SET last_visited=$2 WHERE id=$1;`

// TODO дополнить
func (r *PostgresStorage) CreatePaste(ctx context.Context, tx pgx.Tx, paste domain.Paste, userID int) (uuid.UUID, error) {
	var pasteID uuid.UUID
	row := tx.QueryRow(ctx, createPasteQuery, paste.Title, paste.ExpiresAt, paste.Visibility, userID)
	err := row.Scan(&pasteID)
	return pasteID, err
}

func (r *PostgresStorage) GetPasteByID(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID) (domain.Paste, error) {
	var paste domain.Paste
	row := tx.QueryRow(ctx, getPasteByIDQuery, pasteID)
	err := row.Scan(&paste.Title, &paste.CreatedAt, &paste.ExpiresAt, &paste.Visibility, &paste.LastVisited, &paste.UserID)
	return paste, err
}

func (r *PostgresStorage) UpdateLastVisited(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID) error {
	_, err := tx.Exec(ctx, updateLastVisitedQuery, pasteID, time.Now())
	return err
}

func (r *PostgresStorage) GetUsersPastes(ctx context.Context, userID int) ([]domain.Paste, error) {
	var pastes []domain.Paste
	rows, err := r.db.Pool.Query(ctx, getUsersPastesQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var paste domain.Paste
		if err = rows.Scan(&paste.ID, &paste.Title, &paste.CreatedAt, &paste.ExpiresAt, &paste.Visibility, &paste.LastVisited); err != nil {
			return nil, err
		}
		pastes = append(pastes, paste)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pastes, nil
}

func (r *PostgresStorage) DeletePaste() error {
	return nil
}
