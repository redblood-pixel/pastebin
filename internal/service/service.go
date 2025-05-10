package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/redblood-pixel/pastebin/internal/domain"
	"github.com/redblood-pixel/pastebin/internal/repository"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
	"github.com/redblood-pixel/pastebin/pkg/tokenutil"
)

type Users interface {
	CreateUser(ctx context.Context, name, email, password string) (domain.Tokens, error)
	SignIn(ctx context.Context, nameOrEmail, password_hashed string) (domain.Tokens, error)
	Refresh(ctx context.Context, refresh uuid.UUID) (domain.Tokens, error)
	GetUserById(ctx context.Context, userID int) (domain.User, error)
}

type Pastes interface {
	CreatePaste(ctx context.Context, userID int, paste domain.Paste, data []byte) (string, error)
	GetUsersPastes(ctx context.Context, userID int) ([]domain.Paste, error)
	GetPasteByID(ctx context.Context, pasteID uuid.UUID, userID int, params domain.PasteParameters) (domain.Paste, []byte, error)
	DeletePasteByID(ctx context.Context, pasteID uuid.UUID, userID int) error
}

type Service struct {
	// Set of service interfaces
	Users
	Pastes
}

type Deps struct {
	Postgres     *postgres.Postgres
	TokenManager *tokenutil.TokenManager
	Repository   *repository.Repository
}

func New(deps Deps) *Service {
	return &Service{
		Users:  NewUserService(deps.Postgres, deps.TokenManager, deps.Repository.Database),
		Pastes: NewPastesService(deps.Postgres, deps.Repository),
	}
}
