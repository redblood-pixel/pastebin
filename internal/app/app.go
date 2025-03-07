package app

import (
	"log/slog"

	"github.com/redblood-pixel/pastebin/internal/config"
)

func Run(configPath string) {

	cfg := config.MustLoad(configPath)
	slog.Info(
		"config",
		slog.Any("config", cfg),
	)

	slog.Info("Server started")

	// TODO db conn

	// TODO repository

	// TODO service

	// TODO handler

	// TODO run server

	// TODO graceful shutdown
}
