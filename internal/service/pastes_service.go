package service

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/redblood-pixel/pastebin/internal/domain"
	"github.com/redblood-pixel/pastebin/internal/repository"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
)

type PastesService struct {
	pg *postgres.Postgres
	r  *repository.Repository
}

func NewPastesService(pg *postgres.Postgres, repo *repository.Repository) *PastesService {
	return &PastesService{pg: pg, r: repo}
}

func (s *PastesService) CreatePaste(ctx context.Context, userID int, paste domain.Paste, data []byte) (string, error) {

	conn, err := s.pg.Pool.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	if paste.TTL == 0 {
		paste.TTL = domain.DefaultTTL
	}
	if paste.Visibility == "" {
		paste.Visibility = domain.PublicType
	}
	pasteID, err := s.r.Database.CreatePaste(ctx, tx, paste, userID)
	if err != nil {
		return "", err
	}

	pasteName := getPasteName(userID, pasteID, paste.Title)
	err = s.r.Storage.CreatePaste(ctx, pasteName, paste.TTL, data)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return pasteID.String(), err
}

func (s *PastesService) GetUsersPastes(ctx context.Context, userID int) ([]domain.Paste, error) {

	pastes, err := s.r.GetUsersPastes(ctx, userID)
	return pastes, err
}

func (s *PastesService) GetPasteByID(ctx context.Context, pasteID uuid.UUID, userID int) (domain.Paste, []byte, error) {

	var paste domain.Paste
	conn, err := s.pg.Pool.Acquire(ctx)
	if err != nil {
		return paste, nil, err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return paste, nil, err
	}
	defer tx.Rollback(ctx)

	paste, err = s.r.Database.GetPasteByID(ctx, tx, pasteID)
	if err != nil {
		return paste, nil, err
	}

	// Access control
	if paste.Visibility == domain.PrivateType && paste.UserID != userID {
		return paste, nil, domain.ErrPastePermissionDenied
	}

	err = s.r.Database.UpdateLastVisited(ctx, tx, pasteID)
	if err != nil {
		return paste, nil, err
	}

	pasteName := getPasteName(paste.UserID, pasteID, paste.Title)
	content, err := s.r.Storage.GetPaste(ctx, pasteName)
	// content := []byte(pasteName)
	if err != nil {
		return paste, nil, err
	}
	return paste, content, nil
}

// TODO delete paste

func getPasteName(userID int, pasteID uuid.UUID, title string) string {
	return strconv.Itoa(userID) + "/" + pasteID.String() + "/" + title
}
