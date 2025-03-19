package service

import (
	"context"

	"github.com/redblood-pixel/pastebin/pkg/postgres_queries"
)

type UserSerivce struct {
	db *postgres_queries.Queries
}

func (u *UserSerivce) Create(ctx context.Context, name, email, password string) (int, error) {
	id, err := u.db.CreateUser(ctx, postgres_queries.CreateUserParams{
		Name:           name,
		Email:          email,
		PasswordHashed: password,
	})
	return int(id), err
}
