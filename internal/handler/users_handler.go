package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redblood-pixel/pastebin/pkg/logger"
)

type SignUpInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInInput struct {
	Password    string `json:"password"`
	NameOrEmail string `json:"name_or_email"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

// TODO custom error handling
// TODO Update user data

func (h *Handler) userSignUp(c echo.Context) error {

	var (
		err   error
		input SignUpInput
	)
	logger := logger.WithSource("handler.userSignUp")

	if err = c.Bind(&input); err != nil {
		logger.Error("bind error", "err", err.Error())
		return err
	}

	tokens, err := h.services.Users.CreateUser(
		c.Request().Context(), input.Name, input.Email, input.Password)
	if err != nil {
		logger.Error("Signup error", "err", err.Error())
		return err
	}

	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) userSignIn(c echo.Context) error {
	var (
		err   error
		input SignInInput
	)
	logger := logger.WithSource("handler.userSignIn")

	if err = c.Bind(&input); err != nil {
		logger.Error("binding error", "err", err.Error())
		return err
	}

	tokens, err := h.services.Users.SignIn(c.Request().Context(),
		input.NameOrEmail, input.Password)
	if err != nil {
		logger.Error("Signin error", "err", err.Error())
		return err
	}

	return c.JSON(http.StatusOK, tokens)
}

// TODO make with cookie
func (h *Handler) userRefreshToken(c echo.Context) error {

	var (
		err   error
		input RefreshInput
	)
	logger := logger.WithSource("handler.userRefreshToken")
	if err = c.Bind(&input); err != nil {
		logger.Error("binding error", "err", err.Error())
		return err
	}

	refreshToken, err := uuid.Parse(input.RefreshToken)
	if err != nil {
		logger.Error("not a uuid", "err", err.Error())
		return err
	}
	tokens, err := h.services.Users.Refresh(c.Request().Context(), refreshToken)
	if err != nil {
		logger.Error("refresh error", "err", err.Error())
	}
	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) getUserById(c echo.Context) error {
	return nil
}

func (h *Handler) updateUserById(c echo.Context) error {
	return nil
}
