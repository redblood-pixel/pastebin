package postgres_storage

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

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

const createPastePassword = `INSERT INTO pastes_passwords
(paste_id, password_hashed)
VALUES ($1, $2);`

const getPasteByIDQuery = `SELECT
	title,
	created_at,
	expires_at,
	visibility,
	last_visited,
	burn_after_read,
	user_id,
	pp.password_hashed
FROM pastes AS p
LEFT JOIN pastes_passwords pp ON p.id=pp.paste_id
WHERE p.id=$1
FOR UPDATE OF p;
` // query with locks

const getUsersPastesQuery = `SELECT
	id,
	title,
	created_at,
	expires_at,
	visibility,
	last_visited
FROM pastes
WHERE user_id=$1
`

const updateLastVisitedQuery = `UPDATE pastes SET last_visited=NOW() WHERE id=$1;`

const deletePasteByIDQuery = `DELETE FROM pastes WHERE id=$1;`
const deletePastesQuery = `DELETE FROM pastes WHERE id = ANY($1)`

func (r *PostgresStorage) CreatePaste(ctx context.Context, tx pgx.Tx, paste domain.Paste, userID int) (uuid.UUID, error) {
	var pasteID uuid.UUID
	row := tx.QueryRow(ctx, createPasteQuery, paste.Title, paste.ExpiresAt, paste.Visibility, userID)
	err := row.Scan(&pasteID)
	return pasteID, err
}

func (r *PostgresStorage) CreatePastePassword(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID, passwordHashed string) error {
	_, err := tx.Exec(ctx, createPastePassword, pasteID, passwordHashed)
	return err
}

func (r *PostgresStorage) GetPasteByID(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID) (domain.Paste, error) {
	var paste domain.Paste
	row := tx.QueryRow(ctx, getPasteByIDQuery, pasteID)
	err := row.Scan(
		&paste.Title,
		&paste.CreatedAt,
		&paste.ExpiresAt,
		&paste.Visibility,
		&paste.LastVisited,
		&paste.BurnAfterRead,
		&paste.UserID,
		&paste.Password,
	)
	if err == pgx.ErrNoRows {
		return paste, domain.ErrPasteNotFound
	}
	return paste, err
}

func (r *PostgresStorage) UpdateLastVisited(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID) error {
	_, err := tx.Exec(ctx, updateLastVisitedQuery, pasteID)
	return err
}

func (r *PostgresStorage) GetUsersPastes(ctx context.Context, tx pgx.Tx, userID int, filters domain.PasteFilters) ([]domain.Paste, error) {
	var pastes []domain.Paste
	var query strings.Builder
	var args []any
	paramNum := 2

	query.WriteString(getUsersPastesQuery)
	args = append(args, userID)
	if !filters.CreatedAtFilter.IsZero() {
		query.WriteString(" AND created_at < $" + strconv.Itoa(paramNum))
		args = append(args, filters.CreatedAtFilter)
		paramNum++
	}

	if filters.SortBy != "" {
		query.WriteString(fmt.Sprintf(" ORDER BY %s", filters.SortBy))
		if filters.Desc {
			query.WriteString(" DESC")
		}
	}

	if filters.Limit != 0 {
		query.WriteString(" LIMIT $" + strconv.Itoa(paramNum))
		args = append(args, filters.Limit)
		paramNum++
	}

	if filters.Offset != 0 {
		query.WriteString(" OFFSET $" + strconv.Itoa(paramNum))
		args = append(args, filters.Offset)
		paramNum++
	}
	slog.Debug("getUsersQuery", "query", query.String(), "args", args)

	rows, err := tx.Query(ctx, query.String(), args...)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var paste domain.Paste
		if err = rows.Scan(
			&paste.ID,
			&paste.Title,
			&paste.CreatedAt,
			&paste.ExpiresAt,
			&paste.Visibility,
			&paste.LastVisited,
		); err != nil {
			return nil, err
		}
		pastes = append(pastes, paste)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pastes, nil
}

func (r *PostgresStorage) DeletePasteByID(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePasteByIDQuery, pasteID)
	return err
}

func (r *PostgresStorage) DeletePastes(ctx context.Context, tx pgx.Tx, pastesID []uuid.UUID) error {
	_, err := tx.Exec(ctx, deletePastesQuery, pastesID)
	return err
}
