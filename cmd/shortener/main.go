package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"

	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/db"
	"github.com/korol8484/shortener/internal/app/handlers"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"github.com/korol8484/shortener/internal/app/logger"
	dbstore "github.com/korol8484/shortener/internal/app/storage/db"
	"github.com/korol8484/shortener/internal/app/storage/file"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	userDBStore "github.com/korol8484/shortener/internal/app/user/storage"
)

var (
	// BuildVersion - version
	BuildVersion string = "N/A"
	// BuildDate - date
	BuildDate string = "N/A"
	// BuildCommit - commit
	BuildCommit string = "N/A"
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
	flag.StringVar(&cfg.DBDsn, "d", "", "set postgresql connection string (DSN)")
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
	var store handlers.Store
	var err error
	var pingable handlers.Pingable
	var jwtUserRep middleware.UserAddRepository

	if cfg.DBDsn != "" {
		dbConn, dbErr := db.NewPgDB(cfg)
		if dbErr != nil {
			return dbErr
		}

		pingable = dbConn

		jwtUserRep, err = userDBStore.NewStorage(dbConn)
		if err != nil {
			return err
		}

		store, err = dbstore.NewStorage(dbConn)
		if err != nil {
			return err
		}
	} else if cfg.FileStoragePath != "" {
		store, err = file.NewFileStore(cfg, memory.NewMemStore())
		if err != nil {
			return err
		}

		jwtUserRep = userDBStore.NewMemoryStore()

		defer func(store handlers.Store) {
			_ = store.Close()
		}(store)
	} else {
		jwtUserRep = userDBStore.NewMemoryStore()
		store = memory.NewMemStore()
	}

	if pingable == nil {
		pingable = handlers.NewPingDummy()
	}

	dh, err := handlers.NewDelete(store, log)
	if err != nil {
		return err
	}

	defer dh.Close()

	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	return http.ListenAndServe(
		cfg.Listen,
		handlers.CreateRouter(store, cfg, log, pingable, jwtUserRep, dh),
	)
}
