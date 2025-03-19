package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redblood-pixel/pastebin/pkg/logger"
)

type SignUpInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO add json tags
type SignInInput struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	PasswordHashed string `json:"password_hashed"`
}

// TODO custom error handling
// TODO User Sign in
// TODO Get User by id
// TODO Refresh
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
	return c.JSON(200, "dfdsf")
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
