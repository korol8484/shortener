package main

import (
	"github.com/korol8484/shortener/internal/app"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	store := app.NewMemStore()
	api := app.NewAPI(store)

	mux := http.NewServeMux()
	mux.HandleFunc(`/{id}`, api.HandleRedirect)
	mux.HandleFunc(`/`, api.HandleShort)

	return http.ListenAndServe(`:8080`, mux)
}
