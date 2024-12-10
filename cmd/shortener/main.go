package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/korol8484/shortener/internal/app/grpc/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/db"
	grpcHandler "github.com/korol8484/shortener/internal/app/grpc"
	"github.com/korol8484/shortener/internal/app/handlers"
	"github.com/korol8484/shortener/internal/app/logger"
	dbstore "github.com/korol8484/shortener/internal/app/storage/db"
	"github.com/korol8484/shortener/internal/app/storage/file"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/korol8484/shortener/internal/app/usecase"
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

	if cfg.Grpc {
		if err = runGRPC(cfg, zLog); err != nil {
			zLog.Fatal("can't run application", zap.Error(err))
		}
	} else {
		if err = runHTTP(cfg, zLog); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				zLog.Fatal("can't run application", zap.Error(err))
			}

			zLog.Info("Application shutdown")
		}
	}
}

func runGRPC(cfg *config.App, log *zap.Logger) error {
	//l := strings.SplitN(cfg.Listen, ":", 2)
	//if len(l) != 2 {
	//	return fmt.Errorf("can't parse listen string: %s", cfg.Listen)
	//}

	uCase, jwtUserRep, err := baseInit(cfg, log)
	if err != nil {
		return err
	}
	defer uCase.Close()

	listener, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return err
	}

	gOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(log)),
			grpcHandler.JwtInterceptor(usecase.NewJwt(jwtUserRep, log, "1234567891")),
			grpcHandler.IpInterceptor(cfg.TrustedSubnet, []string{"/service.Internal/Stats"}),
		),
	}

	if cfg.HTTPS.Enable {
		creds, tlsErr := credentials.NewServerTLSFromFile(cfg.HTTPS.Pem, cfg.HTTPS.Key)
		if tlsErr != nil {
			return tlsErr
		}

		gOpts = append(gOpts, grpc.Creds(creds))
	}

	s := grpc.NewServer(gOpts...)
	h := grpcHandler.NewHandler(uCase)
	service.RegisterInternalServer(s, h)

	return s.Serve(listener)
}

func runHTTP(cfg *config.App, log *zap.Logger) error {
	uCase, jwtUserRep, err := baseInit(cfg, log)
	if err != nil {
		return err
	}
	defer uCase.Close()

	stats, err := handlers.NewStats(cfg, log, uCase)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    cfg.Listen,
		Handler: handlers.CreateRouter(uCase, log, jwtUserRep, stats),
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

func baseInit(cfg *config.App, log *zap.Logger) (*usecase.Usecase, usecase.UserAddRepository, error) {
	var store usecase.Store
	var err error
	var pingable usecase.Pingable
	var jwtUserRep usecase.UserAddRepository

	if cfg.DBDsn != "" {
		dbConn, dbErr := db.NewPgDB(cfg)
		if dbErr != nil {
			return nil, nil, dbErr
		}

		pingable = dbConn

		jwtUserRep, err = userDBStore.NewStorage(dbConn)
		if err != nil {
			return nil, nil, err
		}

		store, err = dbstore.NewStorage(dbConn)
		if err != nil {
			return nil, nil, err
		}
	} else if cfg.FileStoragePath != "" {
		store, err = file.NewFileStore(cfg, memory.NewMemStore())
		if err != nil {
			return nil, nil, err
		}

		jwtUserRep = userDBStore.NewMemoryStore()
	} else {
		jwtUserRep = userDBStore.NewMemoryStore()
		store = memory.NewMemStore()
	}

	if pingable == nil {
		pingable = usecase.NewPingDummy()
	}

	uCase := usecase.NewUsecase(
		cfg,
		store,
		pingable,
		log,
	)

	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	return uCase, jwtUserRep, nil
}

func interceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
