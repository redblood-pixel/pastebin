package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/redblood-pixel/pastebin/pkg/logger"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}

func New(ctx context.Context, cfg *Config) (*pgx.Conn, error) {

	var err error
	logger := logger.WithSource("postgres.New")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		logger.Error("error while connecting db", "error", err.Error())
		return nil, err
	}
	if err = conn.Ping(context.Background()); err != nil {
		logger.Error("error while pinging db", "error", err.Error())
		return nil, err
	}
	logger.Info("postgres connected and pinged")

	return conn, nil
}
