package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redblood-pixel/pastebin/internal/service"
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

func applyMiddlewares(router *echo.Echo) {

	// Middleware to limit body size
	router.Use(middleware.BodyLimit("2M"))

	// TODO Maybe implement config for cors
	router.Use(middleware.CORS())

	// Middleware for recovery and human-readable stack trace print
	router.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {

			stackLines := strings.Split(string(stack), "\n")

			for i, line := range stackLines {
				var formattedStack strings.Builder
				if strings.HasPrefix(line, "\t") {
					formattedStack.WriteString("    " + strings.TrimPrefix(line, "\t") + "\n")
				} else {
					formattedStack.WriteString(line + "\n")
				}
				stackLines[i] = formattedStack.String()
			}

			slog.Error(
				"Panic",
				slog.Any("error", err),
				slog.Any("stack", stackLines),
			)
			return FromError(
				service.NewError(
					service.ErrInternalServer,
					errors.New(fmt.Sprintf("panic error: %s", err.Error())),
				),
			)
		},
	}))

	// Logger middleware
	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.String("method", v.Method),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
}
