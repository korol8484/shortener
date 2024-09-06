package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"go.uber.org/zap"
)

func CreateRouter(
	store Store,
	cfg Config,
	logger *zap.Logger,
	p Pingable,
	userRep middleware.UserAddRepository,
) http.Handler {
	api := NewAPI(store, cfg)
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.LoggResponse(logger),
			middleware.LoggRequest(logger),
			middleware.NewCompressor().Handler,
		)

		jwtH := middleware.NewJwt(userRep)

		r.With(jwtH.HandlerSet()).Post("/", api.HandleShort)
		r.Get("/{id}", api.HandleRedirect)
		r.With(jwtH.HandlerSet()).Post("/api/shorten", api.ShortenJSON)
		r.With(jwtH.HandlerSet()).Post("/api/shorten/batch", api.ShortenBatch)
		r.With(jwtH.HandlerRead()).Get("/api/user/urls", api.UserURL)
	})

	r.Get("/ping", Ping(p))

	return r
}
