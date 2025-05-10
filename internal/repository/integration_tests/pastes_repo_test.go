package repository_integration

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redblood-pixel/pastebin/internal/domain"
	"github.com/redblood-pixel/pastebin/internal/repository/postgres_storage"
)

func TestPasteRepository(t *testing.T) {
	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		t.Errorf("connection failed - %s", err.Error())
	}
	defer db.Close(context.Background())
	repo := postgres_storage.NewPostgresStorage(nil)
	fmt.Println("all seted up")
	t.Run("Create", func(t *testing.T) { testCreatePaste(t, repo, db) })
}

func testCreatePaste(t *testing.T, repo *postgres_storage.PostgresStorage, db *pgx.Conn) {

	ctx := context.Background()

	// Подготовка тестовых данных
	now := time.Now()
	testUserID := 1 // Должен существовать в тестовой БД

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
					fmt.Println(err.Error(), reflect.TypeOf(err).String())
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

// func GetUsersPastesTest(t *testing.T) {
// 	db, err := pgx.Connect(context.Background(), dbURL)
// 	if err != nil {
// 		t.Errorf("connection failed - %s", err.Error())
// 	}
// 	defer db.Close(context.Background())

// }
