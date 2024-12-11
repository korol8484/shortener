package handlers

import (
	"context"
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/usecase"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/korol8484/shortener/internal/app/user/storage"
)

func TestDelete_BatchDelete(t *testing.T) {
	router := chi.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()

	j := middleware.NewJwt(usecase.NewJwt(storage.NewMemoryStore(), zap.L(), "123"), zap.L())
	router.Use(j.HandlerSet())

	store := memory.NewMemStore()
	err := store.Add(context.Background(), &domain.URL{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	}, &domain.User{ID: 1})
	require.NoError(t, err)

	defer func(store usecase.Store) {
		_ = store.Close()
	}(store)

	uCase := usecase.NewUsecase(
		&config.App{BaseShortURL: srv.URL},
		store,
		usecase.NewPingDummy(),
		zap.L(),
	)
	defer uCase.Close()

	api := NewAPI(uCase)

	router.Delete("/batch", api.BatchDelete)

	req, err := http.NewRequest("DELETE", srv.URL+"/batch", strings.NewReader("[\"7A2S4z\"]"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	defer req.Body.Close()

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusAccepted)
}
