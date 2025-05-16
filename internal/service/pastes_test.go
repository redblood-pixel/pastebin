//go:build integration
// +build integration

package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/redblood-pixel/pastebin/internal/repository"
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

	// mockStorage := repository.MockStorage{}

	// mock?
	repo := repository.NewRepo(db, nil, "")
	psvc := NewPastesService(repo)

	t.Run("Service Create Paste", func(t *testing.T) { testCreatePaste(t, psvc) })

}

func testCreatePaste(t *testing.T, svc *PastesService) {
	fmt.Println("Hello, world")
}
