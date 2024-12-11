package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"github.com/korol8484/shortener/internal/app/usecase"
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
	uc *usecase.Usecase,
	logger *zap.Logger,
	userRep usecase.UserAddRepository,
	stats *Stats,
) http.Handler {
	api := NewAPI(uc)
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.LoggResponse(logger),
			middleware.LoggRequest(logger),
			middleware.NewCompressor().Handler,
		)

		jwtH := middleware.NewJwt(usecase.NewJwt(userRep, logger, "1234567891"), logger)

		r.With(jwtH.HandlerSet()).Post("/", api.HandleShort)
		r.Get("/{id}", api.HandleRedirect)
		r.With(jwtH.HandlerSet()).Post("/api/shorten", api.ShortenJSON)
		r.With(jwtH.HandlerSet()).Post("/api/shorten/batch", api.ShortenBatch)
		r.With(jwtH.HandlerRead()).Get("/api/user/urls", api.UserURL)
		r.With(jwtH.HandlerRead()).Delete("/api/user/urls", api.BatchDelete)
		r.Get("/api/internal/stats", stats.handle)
	})

	r.Get("/ping", api.Ping)

	return r
}
