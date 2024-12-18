package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/korol8484/shortener/internal/app/user/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_UserURL(t *testing.T) {
	router := chi.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()

	j := middleware.NewJwt(storage.NewMemoryStore(), zap.L(), "123")
	router.Use(j.HandlerSet(), middleware.LoggRequest(zap.L()), middleware.LoggResponse(zap.L()))

	store := memory.NewMemStore()

	defer func(store Store) {
		_ = store.Close()
	}(store)

	err := store.Add(context.Background(), &domain.URL{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	}, &domain.User{ID: 1})
	if err != nil {
		t.Fatal(err)
	}

	api := NewAPI(store, &config.App{BaseShortURL: srv.URL})
	router.Get("/user", api.UserURL)

	req, err := http.NewRequest("GET", srv.URL+"/user", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
