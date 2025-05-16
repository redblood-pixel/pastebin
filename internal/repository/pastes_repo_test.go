//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redblood-pixel/pastebin/internal/domain"
	"github.com/redblood-pixel/pastebin/internal/repository/postgres_storage"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestPasteRepository(t *testing.T) {

	url := os.Getenv("TEST_DB_URL")
	db, err := pgx.Connect(context.Background(), url)
	if err != nil {
		t.Errorf("connection failed - %s", err.Error())
	}
	defer db.Close(context.Background())
	repo := postgres_storage.NewPostgresStorage(nil)
	t.Run("Create", func(t *testing.T) { testCreatePaste(t, repo, db) })
	t.Run("GetUsersPaste", func(t *testing.T) { testGetUsersPaste(t, repo, db) })
}

func testCreatePaste(t *testing.T, repo *postgres_storage.PostgresStorage, db *pgx.Conn) {

	ctx := context.Background()

	now := time.Now()
	testUserID := 1

	tests := []struct {
		name        string
		input       domain.Paste
		userID      int
		wantError   bool
		errorString string // Ожидаемая часть текста ошибки
	}{
		{
			name: "successful creation",
			input: domain.Paste{
				Title:      "Valid Paste",
				ExpiresAt:  now.Add(time.Hour),
				Visibility: "public",
			},
			userID:    testUserID,
			wantError: false,
		},
		{
			name: "empty title",
			input: domain.Paste{
				Title:      "",
				ExpiresAt:  now.Add(time.Hour),
				Visibility: "public",
			},
			userID:    testUserID,
			wantError: false,
		},
		{
			name: "invalid visibility",
			input: domain.Paste{
				Title:      "Invalid Visibility",
				ExpiresAt:  now.Add(time.Hour),
				Visibility: "invalid",
			},
			userID:      testUserID,
			wantError:   true,
			errorString: "invalid input value for enum access_type",
		},
		{
			name: "non-existent user",
			input: domain.Paste{
				Title:      "No Owner",
				ExpiresAt:  now.Add(time.Hour),
				Visibility: "public",
			},
			userID:      99999, // Несуществующий ID
			wantError:   true,
			errorString: " insert or update on table \"pastes\" violates foreign key constraint \"pastes_user_id_fkey\"",
		},
		{
			name: "expired paste",
			input: domain.Paste{
				Title:      "Expired",
				ExpiresAt:  now.Add(-time.Hour), // Прошедшая дата
				Visibility: "public",
			},
			userID:      testUserID,
			wantError:   false,
			errorString: "expiration date must be in future", // ? maybe needs a new constraint
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := db.Begin(ctx)
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback(ctx)

			_, err = repo.CreatePaste(ctx, tx, tt.input, tt.userID)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				} else {
					if tt.errorString != "" && !strings.Contains(err.Error(), tt.errorString) {
						t.Errorf("expected error to contain %q, got %q", tt.errorString, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func testGetUsersPaste(t *testing.T, repo *postgres_storage.PostgresStorage, db *pgx.Conn) {
	// check only titles
	tests := []struct {
		name     string
		userID   int
		filters  domain.PasteFilters
		expected []string
	}{
		{
			name:     "simple get for 1 user",
			userID:   1,
			filters:  domain.PasteFilters{},
			expected: []string{"Public paste 1", "Private paste 1", "Public API docs", "Meeting notes", "Public FAQ"},
		},
		{
			name:     "simple get for 5 user",
			userID:   5,
			filters:  domain.PasteFilters{},
			expected: []string{"x1", "x2", "x3", "x4", "x5", "x6"},
		},
		{
			name:     "simple offset=3 for 5 user",
			userID:   5,
			filters:  domain.PasteFilters{Offset: 3},
			expected: []string{"x4", "x5", "x6"},
		},
		{
			name:     "simple offset=2 limit=2 for 5 user",
			userID:   5,
			filters:  domain.PasteFilters{Offset: 2, Limit: 2},
			expected: []string{"x3", "x4"},
		},
		{
			name:     "simple offset=5 limit=1 for 5 user",
			userID:   5,
			filters:  domain.PasteFilters{Offset: 5, Limit: 1},
			expected: []string{"x6"},
		},
		{
			name:     "created_at_filter=3 hours",
			userID:   5,
			filters:  domain.PasteFilters{CreatedAtFilter: time.Now().Add(-3 * time.Hour)},
			expected: []string{"x2", "x3", "x4", "x5", "x6"},
		},
		{
			name:     "created_at_filter=24 hours",
			userID:   5,
			filters:  domain.PasteFilters{CreatedAtFilter: time.Now().Add(-24 * time.Hour)},
			expected: []string{"x2", "x3", "x4"},
		},
		{
			name:     "created_at_filter=24 hours sort by created_at",
			userID:   5,
			filters:  domain.PasteFilters{CreatedAtFilter: time.Now().Add(-24 * time.Hour), SortBy: "created_at"},
			expected: []string{"x4", "x3", "x2"},
		},
		{
			name:     "created_at_filter=24 hours sort by created_at desc",
			userID:   5,
			filters:  domain.PasteFilters{CreatedAtFilter: time.Now().Add(-24 * time.Hour), SortBy: "created_at", Desc: true},
			expected: []string{"x2", "x3", "x4"},
		},
		{
			name:     "sort by title asc",
			userID:   5,
			filters:  domain.PasteFilters{SortBy: "title"},
			expected: []string{"x1", "x2", "x3", "x4", "x5", "x6"},
		},
		{
			name:     "sort by title desc",
			userID:   5,
			filters:  domain.PasteFilters{SortBy: "title", Desc: true},
			expected: []string{"x6", "x5", "x4", "x3", "x2", "x1"},
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := db.Begin(ctx)
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback(ctx)

			pastes, err := repo.GetUsersPastes(ctx, tx, tt.userID, tt.filters)
			if err != nil {
				t.Errorf("unexpected error - %s", err.Error())
			}
			if len(pastes) != len(tt.expected) {
				titles := make([]string, 0)
				for _, paste := range pastes {
					titles = append(titles, paste.Title)
					fmt.Println(paste.Title, paste.CreatedAt)
				}
				fmt.Println(tt.filters.CreatedAtFilter)
				t.Error(tt.expected)
				t.Error(titles)
				t.Error("expected first, got second")
				t.Fatalf("expected %d rows found %d", len(tt.expected), len(pastes))
			}
			for i := range pastes {
				if pastes[i].Title != tt.expected[i] {
					titles := make([]string, 0)
					for _, paste := range pastes {
						titles = append(titles, paste.Title)
					}
					t.Error(tt.expected)
					t.Error(titles)
					t.Fatal("expected first, got second")
				}
			}
		})
	}
}
