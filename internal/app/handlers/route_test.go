package handlers

import (
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/korol8484/shortener/internal/app/user/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestCreateRouter(t *testing.T) {
	store := memory.NewMemStore()
	cfg := &config.App{}
	pi := &PingDummy{}
	uRep := storage.NewMemoryStore()

	api, err := NewDelete(store, zap.L())
	require.NoError(t, err)
	defer api.Close()

	stat, err := NewStats(cfg, zap.L(), store)
	require.NoError(t, err)

	r := CreateRouter(store, cfg, zap.L(), pi, uRep, api, stat)
	if r == nil {
		t.Fatal("not implement http.Handler")
	}
}
