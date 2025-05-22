package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/redblood-pixel/pastebin/internal/server"
	"github.com/redblood-pixel/pastebin/pkg/minio_connection"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
	"github.com/redblood-pixel/pastebin/pkg/tokenutil"
)

type Config struct {
	Env      string                  `yaml:"env"`
	HTTP     server.Config           `yaml:"http"`
	Postgres postgres.Config         `yaml:"postgres"`
	JWT      tokenutil.Config        `yaml:"jwt"`
	Minio    minio_connection.Config `yaml:"minio"`

	PasswordSalt string `yaml:"password_salt"`
}

func MustLoad(configPath string) *Config {

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Sprintf("error occured while reading config: %s", err.Error()))
	}
	if cfg.Env != "dev" && cfg.Env != "test" && cfg.Env != "prod" {
		panic("error occured while reading config - not valid env value")
	}
	return &cfg
}
