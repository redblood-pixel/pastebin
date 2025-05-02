package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redblood-pixel/pastebin/internal/config"
	"github.com/redblood-pixel/pastebin/internal/handler"
	"github.com/redblood-pixel/pastebin/internal/repository"
	"github.com/redblood-pixel/pastebin/internal/server"
	"github.com/redblood-pixel/pastebin/internal/service"
	"github.com/redblood-pixel/pastebin/pkg/logger"
	"github.com/redblood-pixel/pastebin/pkg/minio_connection"
	"github.com/redblood-pixel/pastebin/pkg/postgres"
	"github.com/redblood-pixel/pastebin/pkg/tokenutil"
)

func Run(configPath string) {

	// Configs, Logger and DB connection
	cfg := config.MustLoad(configPath)
	logger.Init(cfg.Env)

	logger := logger.WithSource("api.Run")

	logger.Debug("Hello word!")

	dbctx := context.Background()
	pg, err := postgres.New(dbctx, &cfg.Postgres)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	minioC, err := minio_connection.Connect()
	if err != nil {
		fmt.Println("min con", err.Error())
		return
	}

	c, err := minioC.BucketExists(dbctx, "pastes")
	if err != nil {
		fmt.Println("min ping", err.Error())
		return
	}
	fmt.Println("bucket -", c)

	// testObj := []byte("print(\"Hello, world!\")")
	// testReader := bytes.NewReader(testObj)
	// uploadInfo, err := minioC.PutObject(dbctx, "pastes", "test_obj", testReader, int64(len(testObj)), minio.PutObjectOptions{})
	// if err != nil {
	// 	fmt.Println("min upload", err.Error())
	// 	return
	// }
	// logger.Info("upload", slog.Any("struct", uploadInfo))
	// Pkg

	// Получение URL для загруженного объекта
	// url, err := minioC.PresignedGetObject(context.Background(), "pastes", "test_obj", time.Second*24*60*60, nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(url)
	tokenManager := tokenutil.New(&cfg.JWT)

	//Service Dependencies
	repository := repository.NewRepo(pg, minioC, "pastes")
	deps := service.Deps{
		Postgres:     pg,
		TokenManager: tokenManager,
		Repository:   repository,
	}

	service := service.New(deps)
	handler := handler.New(service, tokenManager)
	srv := server.New(&cfg.HTTP, handler.Init())

	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Error starting server ", "error", err.Error())
		}
	}()
	logger.Info("Server started")

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
	pg.Close()
	logger.Info("Server stoped")
}
