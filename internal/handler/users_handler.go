package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SignUpInput struct {
	Name     string `json:"name" validate:"required,min=5,max=64"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

type SignInInput struct {
	Password    string `json:"password" validate:"required,min=6,max=64"`
	NameOrEmail string `json:"name_or_email" validate:"required,min=6"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required,uuid4"`
}

// TODO Update user data

func (h *Handler) userSignUp(c echo.Context) error {

	var (
		err   error
		input SignUpInput
	)

	if err = c.Bind(&input); err != nil {
		return err
	}

	err = h.v.Struct(input)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}

	tokens, err := h.services.Users.CreateUser(
		c.Request().Context(), input.Name, input.Email, input.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) userSignIn(c echo.Context) error {
	var (
		err   error
		input SignInInput
	)

	if err = c.Bind(&input); err != nil {
		return err
	}

	err = h.v.Struct(input)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}

	tokens, err := h.services.Users.SignIn(c.Request().Context(),
		input.NameOrEmail, input.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) userRefreshToken(c echo.Context) error {

	var (
		err   error
		input RefreshInput
	)

	if err = c.Bind(&input); err != nil {
		return err
	}

	err = h.v.Struct(input)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}
	refreshToken, err := uuid.Parse(input.RefreshToken)
	if err != nil {
		return err
	}
	tokens, err := h.services.Users.Refresh(c.Request().Context(), refreshToken)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) getUserById(c echo.Context) error {
	var (
		err    error
		userID int
	)

	userIDstr := c.Param("id")
	if userID, err = strconv.Atoi(userIDstr); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.services.Users.GetUserById(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUserById(c echo.Context) error {
	return nil
}
