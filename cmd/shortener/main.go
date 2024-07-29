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
)

func main() {
	cfg := &config.App{}

	flag.StringVar(&cfg.Listen, "a", ":8081", "Http service list addr")
	flag.StringVar(&cfg.BaseShortURL, "b", "http://localhost:8081", "Base short url")
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
	store := storage.NewMemStore()

	return http.ListenAndServe(
		cfg.Listen,
		handlers.CreateRouter(store, cfg, log),
	)
}
