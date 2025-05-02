package postgres_storage

import (
	"github.com/redblood-pixel/pastebin/pkg/postgres"
)

type PostgresStorage struct {
	db *postgres.Postgres
}

func NewPostgresStorage(db *postgres.Postgres) *PostgresStorage {
	return &PostgresStorage{db: db}
}
