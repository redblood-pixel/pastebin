package service

import (
	"context"

	"github.com/redblood-pixel/pastebin/db"
)

type Service struct {
	// Set of service interfaces
	Users
	Pastes
}

type Users interface {
	Create(ctx context.Context, name, email, password string) (int, error)
}

type Pastes interface {
}

func New(querier *db.Queries) *Service {
	return &Service{
		Users: &UserSerivce{
			db: querier,
		},
		Pastes: PastesService{
			db: querier,
		},
	}
}
