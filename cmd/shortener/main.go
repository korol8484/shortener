package main

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	server := &http.Server{
		Addr:    cfg.Listen,
		Handler: handlers.CreateRouter(store, cfg, log, pingable, jwtUserRep, dh),
	}

	oss, stop, errCh := make(chan os.Signal, 1), make(chan struct{}, 1), make(chan error, 1)
	signal.Notify(oss, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-oss

		stop <- struct{}{}
	}()

	if cfg.HTTPS.Enable {
		go func() {
			if err = server.ListenAndServeTLS(
				cfg.HTTPS.Pem,
				cfg.HTTPS.Key,
			); err != nil {
				errCh <- err
			}
		}()
	} else {
		go func() {
			if err = server.ListenAndServe(); err != nil {
				errCh <- err
			}
		}()
	}

	for {
		select {
		case e := <-errCh:
			return e
		case <-stop:
			withTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err = server.Shutdown(withTimeout); err != nil {
				cancel()

				return err
			}

			cancel()

			return nil
		}
	}
}
