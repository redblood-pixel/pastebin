package domain

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrInternalServer = errors.New("internal server error")
	ErrUserExists     = errors.New("user with such email or name already exists")
	ErrRefreshExpired = errors.New("refresh token expired")
)
