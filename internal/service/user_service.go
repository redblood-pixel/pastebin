package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/gommon/log"
	"github.com/redblood-pixel/pastebin/internal/domain"
	"github.com/redblood-pixel/pastebin/pkg/hash"
	"github.com/redblood-pixel/pastebin/pkg/logger"
	"github.com/redblood-pixel/pastebin/pkg/postgres_queries"
	"github.com/redblood-pixel/pastebin/pkg/tokenutil"
)

type UserSerivce struct {
	q    *postgres_queries.Queries
	conn *pgx.Conn

	tm *tokenutil.TokenManager
}

func NewUserService(q *postgres_queries.Queries, conn *pgx.Conn,
	tm *tokenutil.TokenManager) *UserSerivce {
	return &UserSerivce{
		q:    q,
		conn: conn,
		tm:   tm,
	}
}

func (u *UserSerivce) CreateUser(ctx context.Context, name, email, password string) (domain.Tokens, error) {
	var tokens domain.Tokens
	logger := logger.WithSource("service.UserService.CreateUser")

	tx, err := u.conn.Begin(ctx)
	if err != nil {
		logger.Error("transaction is not even begin", "err", err.Error())
		return tokens, err
	}
	defer tx.Rollback(ctx)

	qtx := u.q.WithTx(tx)
	userID, err := qtx.CreateUser(ctx, postgres_queries.CreateUserParams{
		Name:           name,
		Email:          email,
		PasswordHashed: hash.Generate(password),
	})
	if err != nil {
		logger.Error("error while creating user", "err", err.Error())
		logger.Error(err.Error())
		return tokens, err
	}

	// creating tokens
	tokens.AccessToken, err = u.tm.CreateAccessToken(int(userID))
	if err != nil {
		logger.Error("error while creating access token", "err", err.Error())
		return tokens, err
	}
	refresh, err := qtx.CreateSession(ctx, postgres_queries.CreateSessionParams{
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(u.tm.GetRefreshTTL()),
	})
	if err != nil {
		logger.Error("error while creating refresh token", "err", err.Error())
		return tokens, err
	}
	tokens.RefreshToken = refresh.String()

	if err = tx.Commit(ctx); err != nil {
		log.Error("error while commiting tx", "err", err.Error())
		return tokens, err
	}
	return tokens, nil
}

func (u *UserSerivce) SignIn(ctx context.Context, name, email, password_hashed string) (domain.Tokens, error) {
	var tokens domain.Tokens
	logger := logger.WithSource("service.UserService.SignIn")

	tx, err := u.conn.Begin(ctx)
	if err != nil {
		logger.Error("transaction is not even begin", "err", err.Error())
		return tokens, err
	}
	defer tx.Rollback(ctx)

	// qtx := u.q.WithTx(tx)

	return tokens, err
}
