package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redblood-pixel/pastebin/pkg/logger"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
	URL      string
}

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *Config) (*Postgres, error) {

	var (
		err error
		dsn string
	)

	logger := logger.WithSource("postgres.New")
	if cfg.URL == "" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	} else {
		dsn = cfg.URL
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	// Ограничиваем пул
	config.MaxConns = 5
	config.MinConns = 2
	config.MaxConnIdleTime = time.Minute * 5
	config.MaxConnLifetime = time.Hour

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("error while connecting db", "error", err.Error())
		return nil, err
	}
	if err = db.Ping(ctx); err != nil {
		logger.Error("error while pinging db", "error", err.Error())
		return nil, err
	}
	logger.Info("postgres connected and pinged")

	return &Postgres{Pool: db}, nil
}

func (p *Postgres) Close() {
	p.Pool.Close()
}
