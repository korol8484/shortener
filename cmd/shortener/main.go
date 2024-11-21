package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"

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

// Build variables
var (
	// BuildVersion - version
	BuildVersion string = "N/A"
	// BuildDate - date
	BuildDate string = "N/A"
	// BuildCommit - commit
	BuildCommit string = "N/A"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("can't initalize config %s", err)
	}

	zLog, err := logger.NewLogger(false)
	if err != nil {
		log.Fatalf("can't initalize logger %s", err)
	}

	defer func(zLog *zap.Logger) {
		_ = zLog.Sync()
	}(zLog)

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

	if cfg.Https.Enable {
		return http.ListenAndServeTLS(
			cfg.Listen,
			cfg.Https.Pem,
			cfg.Https.Key,
			handlers.CreateRouter(store, cfg, log, pingable, jwtUserRep, dh),
		)
	}

	return http.ListenAndServe(
		cfg.Listen,
		handlers.CreateRouter(store, cfg, log, pingable, jwtUserRep, dh),
	)
}
