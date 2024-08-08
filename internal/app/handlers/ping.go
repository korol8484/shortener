package handlers

import (
	"net/http"
)

type Pingable interface {
	Ping() error
}

func Ping(p Pingable) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := p.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
}
