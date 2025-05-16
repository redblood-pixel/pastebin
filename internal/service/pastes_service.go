package service

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/redblood-pixel/pastebin/internal/domain"
	"github.com/redblood-pixel/pastebin/internal/repository"
	"github.com/redblood-pixel/pastebin/pkg/hash"
)

type PastesService struct {
	r *repository.Repository
}

func NewPastesService(repo *repository.Repository) *PastesService {
	return &PastesService{r: repo}
}

func (s *PastesService) CreatePaste(ctx context.Context, userID int, paste domain.Paste, data []byte) (string, error) {

	tx, err := s.r.CreateTx(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	if paste.ExpiresAt.IsZero() {
		paste.ExpiresAt = time.Now().Add(domain.DefaultTTL)
	}
	if paste.Visibility == "" {
		paste.Visibility = domain.PublicType
	}
	pasteID, err := s.r.Database.CreatePaste(ctx, tx, paste, userID)
	if err != nil {
		return "", err
	}

	if paste.Password.String != "" {
		if err = s.r.Database.CreatePastePassword(ctx, tx, pasteID, hash.Generate9(paste.Password.String)); err != nil {
			return "", err
		}
	}

	pasteName := getPasteName(userID, pasteID, paste.Title)
	err = s.r.Storage.CreatePaste(ctx, pasteName, paste.ExpiresAt, data)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return pasteID.String(), err
}

func (s *PastesService) GetUsersPastes(ctx context.Context, userID int, filters domain.PasteFilters) ([]domain.Paste, error) {

	tx, err := s.r.Database.CreateTx(ctx)
	if err != nil {
		return nil, err
	}
	pastes, err := s.r.GetUsersPastes(ctx, tx, userID, filters)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiredPastes := make([]uuid.UUID, 0)
	resultPastes := make([]domain.Paste, 0)
	for i := range pastes {
		if pastes[i].ExpiresAt.Before(now) {
			expiredPastes = append(expiredPastes, pastes[i].ID)
		} else {
			resultPastes = append(resultPastes, pastes[i])
		}
	}

	if len(expiredPastes) > 0 {
		// Delete expired pastes with go func
		go func() {
			err := s.r.Storage.DeletePastes(context.Background(), userID, expiredPastes)
			if err != nil {
				slog.Debug("can not delete from storage - ",
					slog.String("func", "PastesService.GetUsersPastes"),
					"err", err.Error(),
				)
				tx.Rollback(context.Background())
			}
			err = s.r.Database.DeletePastes(context.Background(), tx, expiredPastes)
			if err != nil {
				slog.Debug("can not delete from db - ",
					slog.String("func", "PastesService.GetUsersPastes"),
					"err", err.Error(),
				)
				tx.Rollback(context.Background())
			}
			if err = tx.Commit(context.Background()); err != nil {
				slog.Debug("can not commmit - ",
					slog.String("func", "PastesService.GetUsersPastes"),
					"err", err.Error(),
				)
			}
			slog.Debug("vacuum in go func ended")
		}()
	}

	return resultPastes, nil
}

func (s *PastesService) GetPasteByID(ctx context.Context, pasteID uuid.UUID, userID int, pastePassword string) (domain.Paste, []byte, error) {

	var paste domain.Paste
	tx, err := s.r.Database.CreateTx(ctx)
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
		if paste.Password.Status == pgtype.Present && hash.CheckPassword(pastePassword, paste.Password.String) {
			goto accessPassed
		}

		return paste, nil, domain.ErrPastePermissionDenied
	}
accessPassed:

	// TTL control
	if paste.ExpiresAt.Before(time.Now()) || // if paste is expired
		paste.LastVisited.Add(domain.DefaultLastVisitedTTL).Before(time.Now()) { // if paste was no visited in 2 years
		err = s.DeletePaste(ctx, tx, pasteID, userID, paste.Title)
		if err != nil {
			return paste, nil, err
		}
		if err = tx.Commit(ctx); err != nil {
			return paste, nil, err
		}
		return paste, nil, domain.ErrPasteExpired
	}

	err = s.r.Database.UpdateLastVisited(ctx, tx, pasteID)
	if err != nil {
		return paste, nil, err
	}

	pasteName := getPasteName(paste.UserID, pasteID, paste.Title)
	content, err := s.r.Storage.GetPaste(ctx, pasteName)
	if err != nil {
		return paste, nil, err
	}

	// Burn after read control
	if paste.BurnAfterRead {
		err = s.DeletePaste(ctx, tx, pasteID, userID, paste.Title)
		if err != nil {
			return paste, content, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return paste, nil, err
	}
	return paste, content, nil
}

func (s *PastesService) DeletePasteByID(ctx context.Context, pasteID uuid.UUID, userID int) error {
	tx, err := s.r.Database.CreateTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	paste, err := s.r.Database.GetPasteByID(ctx, tx, pasteID)
	if err != nil {
		return err
	}

	if paste.UserID != userID {
		return domain.ErrPasteDeleteDenied
	}

	err = s.DeletePaste(ctx, tx, pasteID, userID, paste.Title)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PastesService) DeletePaste(ctx context.Context, tx pgx.Tx, pasteID uuid.UUID, userID int, title string) error {
	err := s.r.Storage.DeletePaste(ctx, getPasteName(userID, pasteID, title))
	if err != nil {
		return err
	}
	err = s.r.Database.DeletePasteByID(ctx, tx, pasteID)
	return err
}

func getPasteName(userID int, pasteID uuid.UUID, title string) string {
	return strconv.Itoa(userID) + "/" + pasteID.String() + "/" + title
}
