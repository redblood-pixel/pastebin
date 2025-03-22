package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redblood-pixel/pastebin/internal/service"
	"github.com/redblood-pixel/pastebin/pkg/tokenutil"
)

type Handler struct {
	services *service.Service
	tm       *tokenutil.TokenManager
}

type APIError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("http status: %d, %s", e.Status, e.Message)
}

func FromError(err error) APIError {
	var (
		apiError APIError
		svc      service.Error
	)
	if errors.As(err, &svc) {
		apiError.Message = svc.AppErr().Error()
		switch svc.SvcErr() {
		case service.ErrUserNotFound:
			apiError.Status = http.StatusNotFound
		case service.ErrUserExists:
			apiError.Status = http.StatusConflict
		case service.ErrRefreshExpired:
			apiError.Status = http.StatusUnauthorized
		case service.ErrInternalServer:
			apiError.Status = http.StatusInternalServerError
		}
	}
	return apiError
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

	router.HTTPErrorHandler = customErrorHandler

	applyMiddlewares(router)
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

func customErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var (
		code    = http.StatusInternalServerError
		message interface{}
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	} else if he, ok := err.(APIError); ok {
		code = he.Status
		message = he
	}
	// TODO maybe send custom http pages
	c.JSON(code, message)
}
