package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/korol8484/shortener/internal/app/handlers/middleware"
)

// CreateRouter create HTTP router, register endpoints using chi
// @title			Shortener API
// @version		1.0
// @description	Golang service
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	https://localhost:8080
// @contact.email	info@localhost.ru
//
// @host			http://localhost:8080
// @BasePath		/
func CreateRouter(
	store Store,
	cfg Config,
	logger *zap.Logger,
	p Pingable,
	userRep middleware.UserAddRepository,
	deleteHandler *Delete,
) http.Handler {
	api := NewAPI(store, cfg)
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.LoggResponse(logger),
			middleware.LoggRequest(logger),
			middleware.NewCompressor().Handler,
		)

		jwtH := middleware.NewJwt(userRep, logger, "12345dsdsdtoken")

		r.With(jwtH.HandlerSet()).Post("/", api.HandleShort)
		r.Get("/{id}", api.HandleRedirect)
		r.With(jwtH.HandlerSet()).Post("/api/shorten", api.ShortenJSON)
		r.With(jwtH.HandlerSet()).Post("/api/shorten/batch", api.ShortenBatch)
		r.With(jwtH.HandlerRead()).Get("/api/user/urls", api.UserURL)
		r.With(jwtH.HandlerRead()).Delete("/api/user/urls", deleteHandler.BatchDelete)
	})

	r.Get("/ping", Ping(p))

	return r
}
