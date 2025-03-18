package service

import (
	"context"

	"github.com/redblood-pixel/pastebin/db"
)

type UserSerivce struct {
	db *db.Queries
}

func (u *UserSerivce) Create(ctx context.Context, name, email, password string) (int, error) {
	id, err := u.db.CreateUser(ctx, db.CreateUserParams{
		Name:           name,
		Email:          email,
		PasswordHashed: password,
	})
	return int(id), err
}
