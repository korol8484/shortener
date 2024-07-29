package main

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v11"
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/handlers"
	"github.com/korol8484/shortener/internal/app/logger"
	"github.com/korol8484/shortener/internal/app/storage"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	cfg := &config.App{}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't retrive pwd %s", err)
	}

	flag.StringVar(&cfg.Listen, "a", ":8080", "Http service list addr")
	flag.StringVar(&cfg.BaseShortURL, "b", "http://localhost:8080", "Base short url")
	flag.StringVar(&cfg.FileStoragePath, "f", path.Join(pwd, "/data/db"), "set db file path")
	flag.Parse()

	zLog, err := logger.NewLogger(false)
	if err != nil {
		log.Fatalf("can't initalize logger %s", err)
	}

	defer func(zLog *zap.Logger) {
		_ = zLog.Sync()
	}(zLog)

	if err = env.Parse(cfg); err != nil {
		zLog.Warn("can't parse environment variables", zap.Error(err))
	}

	if err = run(cfg, zLog); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			zLog.Fatal("can't run application", zap.Error(err))
		}

		zLog.Info("Application shutdown")
	}
}

func run(cfg *config.App, log *zap.Logger) error {
	store, err := storage.NewMemStore(cfg)
	if err != nil {
		return err
	}

	defer func(store storage.Store) {
		_ = store.Close()
	}(store)

	return http.ListenAndServe(
		cfg.Listen,
		handlers.CreateRouter(store, cfg, log),
	)
}
