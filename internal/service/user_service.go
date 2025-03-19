package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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
		// TODO handle this error
		logger.Error("error while creating user", "err", err.Error())
		logger.Error(err.Error())
		return tokens, err
	}

	tokens, err = u.CreateNewSession(ctx, qtx, userID)

	if err = tx.Commit(ctx); err != nil {
		log.Error("error while commiting tx", "err", err.Error())
		return tokens, err
	}
	return tokens, nil
}

func (u *UserSerivce) SignIn(ctx context.Context, nameOrEmail, password string) (domain.Tokens, error) {
	var tokens domain.Tokens
	logger := logger.WithSource("service.UserService.SignIn")

	tx, err := u.conn.Begin(ctx)
	if err != nil {
		logger.Error("transaction is not even begin", "err", err.Error())
		return tokens, err
	}
	defer tx.Rollback(ctx)

	qtx := u.q.WithTx(tx)
	row, err := qtx.FindUserByNameOrEmail(ctx, nameOrEmail)
	if !hash.CheckPassword(password, row.PasswordHashed) {
		return tokens, domain.ErrNotFound
	}
	if err != nil {
		logger.Error("error user searching in db", "err", err.Error())
		return tokens, err
	}

	tokens, err = u.CreateNewSession(ctx, qtx, row.ID)

	if err = tx.Commit(ctx); err != nil {
		logger.Error("error while commiting tx", "err", err.Error())
	}
	return tokens, err
}

func (u *UserSerivce) GetUserById(ctx context.Context, userID int) (domain.User, error) {
	var user domain.User
	row, err := u.q.GetUserById(ctx, int32(userID))
	if err != nil {
		return user, domain.ErrNotFound
	}
	user.Name = row.Name
	user.CreatedAt = row.CreatedAt
	user.LastLogin = row.LastLogin
	return user, nil
}

func (u *UserSerivce) Refresh(ctx context.Context, refresh uuid.UUID) (domain.Tokens, error) {
	var tokens domain.Tokens
	logger := logger.WithSource("service.UserService.Refresh")

	tx, err := u.conn.Begin(ctx)
	if err != nil {
		logger.Error("transaction is not even begin", "err", err.Error())
	}
	defer tx.Rollback(ctx)

	qtx := u.q.WithTx(tx)
	row, err := qtx.DeleteById(ctx, refresh)
	if time.Now().After(row.ExpiresAt) {
		return tokens, fmt.Errorf("refresh expired")
	}

	tokens, err = u.CreateNewSession(ctx, qtx, row.UserID)
	if err != nil {
		logger.Debug("sessions creation error")
	}
	return tokens, err
}

func (u *UserSerivce) CreateNewSession(ctx context.Context, qtx *postgres_queries.Queries, userID int32) (domain.Tokens, error) {
	var (
		err    error
		tokens domain.Tokens
	)
	logger := logger.WithSource("service.UserService.CreateNewSession")

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
	return tokens, err
}
