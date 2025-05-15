package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/redblood-pixel/pastebin/internal/domain"
)

type CreatePasteInput struct {
	PasteTitle      string `json:"title" validate:"required"`
	PasteTTL        string `json:"ttl"`
	PasteVisibility string `json:"visibility" validate:"omitempty,oneof=public private"`
	PasteContent    string `json:"content" validate:"required,min=8"`
	PastePassword   string `json:"password" validate:"omitempty,min=4,max=16"`
}

type CreatePasteResponse struct {
	PasteID string `json:"paste_id"`
}

type GetUsersPastesInput struct {
	Duration      string `json:"duration"`
	Offset        int    `json:"offset"`
	Limit         int    `json:"limit"`
	SortParameter string `json:"sort_param"`
	Desc          bool   `json:"desc"`
}

type GetPasteInput struct {
	PastePassword string `json:"password"`
}

type GetPasteResponse struct {
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Visibility  string    `json:"visibility"`
	LastVisited time.Time `json:"last_visited"`
	Content     string    `json:"content"`
}

func (h *Handler) createPaste(c echo.Context) error {
	userID := c.Get("userID").(int)
	fmt.Println(userID)

	var (
		err   error
		input CreatePasteInput
	)

	if err = c.Bind(&input); err != nil {
		return err
	}
	err = h.v.Struct(input)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}
	var expiresAt time.Time
	if input.PasteTTL != "" {
		ttl, err := time.ParseDuration(input.PasteTTL)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "ttl must be in time.Duration format")
		}
		if ttl < 0 {
			return c.JSON(http.StatusBadRequest, "ttl must be positive duration")
		}
		expiresAt = time.Now().Add(ttl)
	}
	fmt.Println(expiresAt)

	paste := domain.Paste{
		Title:      input.PasteTitle,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
		Visibility: input.PasteVisibility,
		Password: pgtype.Text{
			String: input.PastePassword,
			Status: pgtype.Present,
		},
	}

	pasteID, err := h.services.Pastes.CreatePaste(c.Request().Context(), userID, paste, []byte(input.PasteContent))
	if err != nil {
		return err
	}
	fmt.Println(pasteID)

	return c.JSON(http.StatusOK, CreatePasteResponse{pasteID})
}

func (h *Handler) getPaste(c echo.Context) error {

	var (
		userID  int
		pasteID uuid.UUID
		input   GetPasteInput
		err     error
	)
	userID, _ = c.Get("userID").(int)
	pasteIDstr := c.Param("id")
	pasteID, err = uuid.Parse(pasteIDstr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "id should be valid uuid")
	}

	if err = c.Bind(&input); err != nil {
		return err
	}

	paste, content, err := h.services.GetPasteByID(c.Request().Context(), pasteID, userID, input.PastePassword)
	if err != nil {
		return err
	}
	response := GetPasteResponse{
		Title:       paste.Title,
		CreatedAt:   paste.CreatedAt,
		ExpiresAt:   paste.ExpiresAt,
		Visibility:  paste.Visibility,
		LastVisited: paste.LastVisited,
		Content:     string(content),
	}
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) getUsersPastes(c echo.Context) error {
	userID := c.Get("userID").(int)
	var (
		input GetUsersPastesInput
		err   error
	)

	if err = c.Bind(&input); err != nil {
		return err
	}

	var createdAtFilter time.Time
	if input.Duration != "" {
		duration, err := time.ParseDuration(input.Duration)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "duration must be a valid duration, i.e. 2m")
		}
		if duration < 0 {
			return c.JSON(http.StatusBadRequest, "ttl must be positive duration")
		}
		createdAtFilter = time.Now().Add(-duration)
	}

	filters := domain.PasteFilters{
		CreatedAtFilter: createdAtFilter,
		SortBy:          input.SortParameter,
		Desc:            input.Desc,
		Limit:           input.Limit,
		Offset:          input.Offset,
	}
	pastes, err := h.services.Pastes.GetUsersPastes(c.Request().Context(), userID, filters)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pastes)
}

func (h *Handler) deletePaste(c echo.Context) error {
	fmt.Println("hel")
	userID, _ := c.Get("userID").(int)
	pasteIDstr := c.Param("id")
	pasteID, err := uuid.Parse(pasteIDstr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "id should be valid uuid")
	}
	err = h.services.Pastes.DeletePasteByID(c.Request().Context(), pasteID, userID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
