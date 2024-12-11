package handlers

import (
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/korol8484/shortener/internal/app/usecase"
	"github.com/korol8484/shortener/internal/app/user/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCreateRouter(t *testing.T) {
	store := memory.NewMemStore()
	cfg := &config.App{}
	pi := &usecase.PingDummy{}
	uRep := storage.NewMemoryStore()

	uCase := usecase.NewUsecase(
		cfg,
		store,
		pi,
		zap.L(),
	)
	defer uCase.Close()

	stat, err := NewStats(cfg, zap.L(), uCase)
	require.NoError(t, err)

	r := CreateRouter(uCase, zap.L(), uRep, stat)
	if r == nil {
		t.Fatal("not implement http.Handler")
	}
}
