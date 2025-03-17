package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/redblood-pixel/pastebin/internal/service"
)

type Handler struct {
	services *service.Service
}

func New(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() *echo.Echo {

	router := echo.New()

	router.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, struct {
			status string
		}{status: "ok"})
	})
	return router
}
