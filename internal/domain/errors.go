package domain

import (
	"errors"
	"net/http"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInternalServer  = errors.New("internal server error")
	ErrSessionNotFound = errors.New("session not found")
	ErrRefreshExpired  = errors.New("refresh token expired")

	ErrPasteNotFound         = errors.New("paste not found")
	ErrPasteExpired          = errors.New("paste has been expired")
	ErrPasteDeleteDenied     = errors.New("paste deletion denied")
	ErrPastePermissionDenied = errors.New("you have no permissions to access paste")
)

var HTTPErrors = map[error]int{
	ErrUserNotFound:    http.StatusNotFound,
	ErrSessionNotFound: http.StatusUnauthorized,
	ErrRefreshExpired:  http.StatusUnauthorized,

	ErrPasteNotFound:         http.StatusNotFound,
	ErrPasteExpired:          http.StatusNotFound,
	ErrPasteDeleteDenied:     http.StatusForbidden,
	ErrPastePermissionDenied: http.StatusForbidden,
	ErrInternalServer:        http.StatusInternalServerError,
}
