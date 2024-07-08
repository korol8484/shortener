package main

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/handlers"
	"github.com/korol8484/shortener/internal/app/storage"
	"net/http"
)

func main() {
	cfg := &config.App{}

	flag.StringVar(&cfg.Listen, "a", ":8080", "Http service list addr")
	flag.StringVar(&cfg.BaseShortURL, "b", "http://localhost:8080", "Base short url")
	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	if err := run(cfg); err != nil {
		panic(err)
	}
}

func run(cfg *config.App) error {
	store := storage.NewMemStore()
	api := handlers.NewAPI(store, cfg)

	r := chi.NewRouter()
	r.Post("/", api.HandleShort)
	r.Get("/{id}", api.HandleRedirect)

	return http.ListenAndServe(cfg.Listen, r)
}
