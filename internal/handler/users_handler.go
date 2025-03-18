package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/redblood-pixel/pastebin/pkg/logger"
)

// TODO User Sign up
// TODO User Sign in
// TODO Get User by id
// TODO Refresh
// TODO Update user data

func (h *Handler) userSignUp(c echo.Context) error {
	logger := logger.WithSource("Handler.userSignUp")
	id, err := h.services.Users.Create(c.Request().Context(), "admin", "ya.ru", "123")
	if err != nil {
		logger.Error("Signup", err.Error())
	}
	logger.Info("id", id)
	return nil
}

func (h *Handler) userSignIn(c echo.Context) error {
	return nil
}

func (h *Handler) userRefreshToken(c echo.Context) error {
	return nil
}

func (h *Handler) getUserById(c echo.Context) error {
	return nil
}

func (h *Handler) updateUserById(c echo.Context) error {
	return nil
}
