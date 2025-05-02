package domain

import (
	"errors"
	"net/http"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInternalServer        = errors.New("internal server error")
	ErrSessionNotFound       = errors.New("session not found")
	ErrRefreshExpired        = errors.New("refresh token expired")
	ErrPastePermissionDenied = errors.New("you have no permissions to access paste")
)

var HTTPErrors = map[error]int{
	ErrUserNotFound:    http.StatusNotFound,
	ErrSessionNotFound: http.StatusUnauthorized,
	ErrRefreshExpired:  http.StatusUnauthorized,
	ErrInternalServer:  http.StatusInternalServerError,
}
