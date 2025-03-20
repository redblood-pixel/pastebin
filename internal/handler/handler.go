package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/redblood-pixel/pastebin/internal/service"
	"github.com/redblood-pixel/pastebin/pkg/tokenutil"
)

type Handler struct {
	services *service.Service
	tm       *tokenutil.TokenManager
}

func New(services *service.Service, tm *tokenutil.TokenManager) *Handler {
	return &Handler{
		services: services,
		tm:       tm,
	}
}

func (h *Handler) Init() *echo.Echo {

	router := echo.New()

	router.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, struct {
			Status string
		}{Status: "ok"})
	})

	h.initRoutes(router)

	return router
}

func (h *Handler) initRoutes(router *echo.Echo) {
	api := router.Group("/api/v1")

	users := api.Group("/users")
	users.POST("/signup", h.userSignUp)
	users.POST("/signin", h.userSignIn)
	users.POST("/refresh", h.userRefreshToken)
	users.GET("/:id", h.getUserById)
	users.PUT("/:id", h.updateUserById)

	pastes := api.Group("/pastes")
	pastes.Use(h.AuthMiddleware)
	pastes.POST("/", h.createPase)
	pastes.GET("/:id", h.getPaste)
	pastes.DELETE("/:id", h.deletePaste)
}
