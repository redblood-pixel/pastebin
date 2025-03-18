package server

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	httpServer http.Server
}

type Config struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

func New(cfg *Config, h http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:         ":" + strconv.Itoa(cfg.Port),
			Handler:      h,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
