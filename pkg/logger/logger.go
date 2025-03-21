package logger

import (
	"log/slog"
	"os"

	"github.com/redblood-pixel/pastebin/pkg/logger/prettylog"
)

// Init инициализирует логгер в зависимости от окружения
func Init(env string) {
	var logger *slog.Logger

	if env != "prod" {
		logger = slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	// Устанавливаем обертку вокруг логгера
	slog.SetDefault(logger)
}

// withSource создает локальный логгер с указанным источником (src)
func WithSource(src string) *slog.Logger {
	return slog.Default().With("src", src)
}
