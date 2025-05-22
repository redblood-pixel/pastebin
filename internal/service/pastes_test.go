//go:build integration
// +build integration

package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/redblood-pixel/pastebin/internal/repository"
	"github.com/redblood-pixel/pastebin/mocks"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
)

func TestPasteService(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	url := os.Getenv("TEST_DB_URL")
	db, err := postgres.New(ctx, &postgres.Config{
		URL: url,
	})

	if err != nil {
		t.Fatalf("can not connect to db - %s", err.Error())
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	// mockStorage.EXPECT().
	// 	GetPaste(gomock.Any()). // Принимаем любой UUID
	// 	DoAndReturn(func(id string) (string, error) {
	// 		if err := uuid.Validate(id); err != nil { // Проверяем, что это валидный UUID
	// 			return "", errors.New("invalid id")
	// 		}
	// 		return "Пример текста для paste", nil
	// 	})

	repo := repository.NewRepo(db, nil, "")
	repo.Storage = mockStorage
	psvc := NewPastesService(repo)

	t.Run("Service Create Paste", func(t *testing.T) { testCreatePaste(t, psvc) })

}

func testCreatePaste(t *testing.T, svc *PastesService) {
	fmt.Println("Hello, world")
}
