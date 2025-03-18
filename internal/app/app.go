package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redblood-pixel/pastebin/db"
	"github.com/redblood-pixel/pastebin/internal/config"
	"github.com/redblood-pixel/pastebin/internal/handler"
	"github.com/redblood-pixel/pastebin/internal/server"
	"github.com/redblood-pixel/pastebin/internal/service"
	logger "github.com/redblood-pixel/pastebin/pkg/logger"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
)

func Run(configPath string) {

	// Configs, Logger and DB connection
	cfg := config.MustLoad(configPath)

	logger.Init(cfg.Env)

	logger := logger.WithSource("api.Run")

	logger.Info(
		"config",
		slog.Any("config", cfg),
	)
	logger.Info("Server started")

	dbctx := context.Background()
	conn, err := postgres.New(dbctx, &cfg.Postgres)
	if err != nil {
		return
	}

	// Generated DB querier, service, handler and server
	querier := db.New(conn)
	service := service.New(querier)
	handler := handler.New(service)
	srv := server.New(&cfg.HTTP, handler.Init())

	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Error starting server ", "error", err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		logger.Error("Error stoping server", "error", err.Error())
	}
	if err := conn.Close(dbctx); err != nil {
		logger.Error("Error closing postgres connection:", "error", err.Error())
	}
	logger.Info("Server stoped")
}
