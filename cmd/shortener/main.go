package main

import (
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

	mux := http.NewServeMux()
	mux.HandleFunc(`/{id}`, api.HandleRedirect)
	mux.HandleFunc(`/`, api.HandleShort)

	return http.ListenAndServe(`:8080`, mux)
}
