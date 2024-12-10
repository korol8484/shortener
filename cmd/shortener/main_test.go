package main

import (
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"os"
	"path"
	"syscall"
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	os.Clearenv()

	go func() {
		time.Sleep(100 * time.Millisecond)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	err := runHTTP(&config.App{HTTPS: &config.HTTPS{}, Listen: ":8099"}, zap.L())
	require.NoError(t, err)

	err = runHTTP(&config.App{HTTPS: &config.HTTPS{}, Listen: "8099"}, zap.L())
	require.Error(t, err)

	go func() {
		time.Sleep(500 * time.Millisecond)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	_ = os.Setenv("SERVER_ADDRESS", ":8098")
	_ = os.Setenv("ENABLE_HTTPS", "1")
	_ = os.Setenv("FILE_STORAGE_PATH", path.Join(os.TempDir(), "/db"))

	cfg, err := config.NewConfig()
	require.NoError(t, err)

	err = runHTTP(cfg, zap.L())
	require.NoError(t, err)

	_ = os.Remove(cfg.HTTPS.Key)
	_ = os.Remove(cfg.HTTPS.Pem)
}
