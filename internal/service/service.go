package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/redblood-pixel/pastebin/internal/domain"
	"github.com/redblood-pixel/pastebin/pkg/postgres_queries"
	"github.com/redblood-pixel/pastebin/pkg/tokenutil"
)

type Users interface {
	CreateUser(ctx context.Context, name, email, password string) (domain.Tokens, error)
	SignIn(ctx context.Context, nameOrEmail, password_hashed string) (domain.Tokens, error)
	Refresh(ctx context.Context, refresh uuid.UUID) (domain.Tokens, error)
}

type Pastes interface {
}

type Service struct {
	// Set of service interfaces
	Users
	Pastes
}

type Deps struct {
	Querier      *postgres_queries.Queries
	PostgresConn *pgx.Conn
	TokenManager *tokenutil.TokenManager
}

func New(deps Deps) *Service {
	return &Service{
		Users: NewUserService(deps.Querier, deps.PostgresConn, deps.TokenManager),
		Pastes: PastesService{
			db: deps.Querier,
		},
	}
}
