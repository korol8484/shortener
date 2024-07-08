package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/handlers"
	"github.com/korol8484/shortener/internal/app/storage"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	store := storage.NewMemStore()
	api := handlers.NewAPI(store)

	r := chi.NewRouter()
	r.Post("/", api.HandleShort)
	r.Get("/{id}", api.HandleRedirect)

	return http.ListenAndServe(`:8080`, r)
}
