package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/redblood-pixel/pastebin/internal/server"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
)

type Config struct {
	Env      string          `yaml:"env"`
	HTTP     server.Config   `yaml:"http"`
	Postgres postgres.Config `yaml:"postgres"`
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
