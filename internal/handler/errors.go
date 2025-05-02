package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redblood-pixel/pastebin/internal/domain"
)

func customErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var (
		code    = http.StatusInternalServerError
		message struct {
			Message interface{} `json:"message"`
		}
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message.Message = he.Message
	} else if he, ok := err.(validator.ValidationErrors); ok {
		code = http.StatusBadRequest
		message.Message = he.Error()
	} else if he, ok := domain.HTTPErrors[err]; ok {
		code = he
		message.Message = err.Error()
	}
	// TODO maybe send custom http pages
	c.JSON(code, message)
}
