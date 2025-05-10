package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/redblood-pixel/pastebin/internal/domain"
	"github.com/redblood-pixel/pastebin/internal/repository/minio_storage"
	"github.com/redblood-pixel/pastebin/internal/repository/postgres_storage"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
)

type Database interface {
	CreateUser(ctx context.Context, tx pgx.Tx, name, email, password_hashed string) (int, error)
	CreateSession(ctx context.Context, tx pgx.Tx, userID int, expireAt time.Duration) (uuid.UUID, error)
	FindUserByNameOrEmail(ctx context.Context, nameOrEmail string) (int, string, error)
	GetUserById(ctx context.Context, userID int) (domain.User, error)
	DeleteSessionById(ctx context.Context, tx pgx.Tx, sessionID uuid.UUID) (int, time.Time, error)

	CreatePaste(ctx context.Context, tx pgx.Tx, paste domain.Paste, userID int) (uuid.UUID, error)
	GetPasteByID(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID) (domain.Paste, error)
	GetUsersPastes(ctx context.Context, userID int, createdAtFilter time.Time, sortBy string, offset, limit int) ([]domain.Paste, error)
	UpdateLastVisited(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID) error
	DeletePasteByID(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID) error
}

type Storage interface {
	CreatePaste(ctx context.Context, name string, ttl time.Time, data []byte) error
	GetPaste(ctx context.Context, name string) ([]byte, error)
	DeletePaste(ctx context.Context, name string) error
}

type Repository struct {
	Database
	Storage
}

func NewRepo(db *postgres.Postgres, mc *minio.Client, bucketName string) *Repository {
	return &Repository{
		Database: postgres_storage.NewPostgresStorage(db),
		Storage:  minio_storage.NewPastesRepository(mc, bucketName),
	}
}
