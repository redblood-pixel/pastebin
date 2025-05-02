package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redblood-pixel/pastebin/internal/domain"
)

// TODO Create Paste
// TODO Get Paste
// TODO Delete Paste

type CreatePasteInput struct {
	PasteTitle      string        `json:"title" validate:"required"`
	PasteTTL        time.Duration `json:"ttl"`
	PasteVisibility string        `json:"visibility_access_type" validate:"omitempty,oneof='public private'"`
	PasteContent    string        `json:"content" validate:"required,min=8"`
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

	paste := domain.Paste{
		Title:      input.PasteTitle,
		CreatedAt:  time.Now(),
		TTL:        input.PasteTTL,
		Visibility: input.PasteVisibility,
	}

	pasteID, err := h.services.Pastes.CreatePaste(c.Request().Context(), userID, paste, []byte(input.PasteContent))
	if err != nil {
		return err
	}
	// Service call
	fmt.Println(pasteID)

	return c.NoContent(http.StatusOK)
}

func (h *Handler) getPaste(c echo.Context) error {
	userID, _ := c.Get("userID").(int)
	pasteIDstr := c.Param("id")
	pasteID, err := uuid.Parse(pasteIDstr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "id should be valid uuid")
	}

	paste, content, err := h.services.GetPasteByID(c.Request().Context(), pasteID, userID)
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
	pastes, err := h.services.Pastes.GetUsersPastes(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pastes)
}

func (h *Handler) deletePaste(c echo.Context) error {
	return nil
}
