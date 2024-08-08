package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"go.uber.org/zap"
	"net/http"
)

func CreateRouter(
	store Store,
	cfg Config,
	logger *zap.Logger,
	p Pingable,
) http.Handler {
	api := NewAPI(store, cfg)
	r := chi.NewRouter()

	r.Use(
		middleware.LoggResponse(logger),
		middleware.LoggRequest(logger),
		middleware.NewCompressor().Handler,
	)

	r.Post("/", api.HandleShort)
	r.Get("/{id}", api.HandleRedirect)
	r.Post("/api/shorten", api.ShortenJSON)
	r.Get("/ping", Ping(p))

	return r
}
