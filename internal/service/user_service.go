package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/redblood-pixel/pastebin/internal/domain"
	"github.com/redblood-pixel/pastebin/internal/repository"
	"github.com/redblood-pixel/pastebin/pkg/hash"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
	"github.com/redblood-pixel/pastebin/pkg/tokenutil"
)

type UserSerivce struct {
	pg *postgres.Postgres
	ur repository.Database
	tm *tokenutil.TokenManager
}

func NewUserService(
	pg *postgres.Postgres,
	tm *tokenutil.TokenManager,
	ur repository.Database,
) *UserSerivce {
	return &UserSerivce{
		pg: pg,
		ur: ur,
		tm: tm,
	}
}

// TODO дописать кастомизацию ошибок

func (u *UserSerivce) CreateUser(ctx context.Context, name, email, password string) (domain.Tokens, error) {
	var tokens domain.Tokens

	conn, err := u.pg.Pool.Acquire(ctx)
	if err != nil {
		return tokens, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return tokens, err
	}
	defer tx.Rollback(ctx)

	userID, err := u.ur.CreateUser(ctx, tx, name, email, hash.Generate(password))
	if err != nil {
		return tokens, err
	}

	tokens, err = u.CreateSession(ctx, tx, userID)
	if err = tx.Commit(ctx); err != nil {
		return tokens, err
	}
	return tokens, nil
}

func (u *UserSerivce) SignIn(ctx context.Context, nameOrEmail, password string) (domain.Tokens, error) {
	var tokens domain.Tokens

	conn, err := u.pg.Pool.Acquire(ctx)
	if err != nil {
		return tokens, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return tokens, err
	}
	defer tx.Rollback(ctx)

	userID, passwordHashed, err := u.ur.FindUserByNameOrEmail(ctx, nameOrEmail)
	if err != nil || !hash.CheckPassword(password, passwordHashed) {
		return tokens, domain.ErrUserNotFound
	}

	tokens, err = u.CreateSession(ctx, tx, userID)
	if err != nil {
		return tokens, err
	}

	if err = tx.Commit(ctx); err != nil {
		return tokens, err
	}

	return tokens, err
}

func (u *UserSerivce) GetUserById(ctx context.Context, userID int) (domain.User, error) {
	user, err := u.ur.GetUserById(ctx, userID)
	return user, err
}

func (u *UserSerivce) Refresh(ctx context.Context, refresh uuid.UUID) (domain.Tokens, error) {
	var tokens domain.Tokens

	conn, err := u.pg.Pool.Acquire(ctx)
	if err != nil {
		return tokens, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return tokens, err
	}
	defer tx.Rollback(ctx)

	userID, expireAt, err := u.ur.DeleteSessionById(ctx, tx, refresh)
	if err != nil {
		return tokens, domain.ErrSessionNotFound
	}
	if time.Now().After(expireAt) {
		return tokens, domain.ErrRefreshExpired
	}

	tokens, err = u.CreateSession(ctx, tx, userID)
	if err != nil {
		return tokens, err
	}

	if err = tx.Commit(ctx); err != nil {
		return tokens, err
	}

	return tokens, err
}

func (u *UserSerivce) CreateSession(ctx context.Context, tx pgx.Tx, userID int) (domain.Tokens, error) {

	var tokens domain.Tokens

	refresh, err := u.ur.CreateSession(ctx, tx, userID, u.tm.GetRefreshTTL())
	if err != nil {
		return tokens, err
	}
	tokens.RefreshToken = refresh.String()
	tokens.AccessToken, err = u.tm.CreateAccessToken(userID)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}
