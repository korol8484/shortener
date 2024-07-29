package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"github.com/korol8484/shortener/internal/app/storage"
	"go.uber.org/zap"
	"net/http"
)

func CreateRouter(
	store storage.Store,
	cfg Config,
	logger *zap.Logger,
) http.Handler {
	api := NewAPI(store, cfg)
	r := chi.NewRouter()

	r.Use(
		middleware.LoggResponse(logger),
		middleware.LoggRequest(logger),
	)

	r.Post("/", api.HandleShort)
	r.Get("/{id}", api.HandleRedirect)
	r.Post("/api/shorten", api.ShortenJson)

	return r
}
