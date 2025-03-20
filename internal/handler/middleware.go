package handler

import (
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return next(c)
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return next(c)
		}

		userID, err := h.tm.ParseAccessToken(headerParts[1])
		if err != nil {
			return next(c)
		}
		c.Set("userID", userID)
		slog.Info("middleware", "id", userID)
		return next(c)
	}
}
