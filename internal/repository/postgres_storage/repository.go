package postgres_storage

import (
	"context"

	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
)

type PostgresStorage struct {
	db *postgres.Postgres
}

func NewPostgresStorage(db *postgres.Postgres) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (r *PostgresStorage) CreateTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		slog.Error("can not create tx", "err", err.Error())
		return nil, err
	}
	return tx, nil
}
