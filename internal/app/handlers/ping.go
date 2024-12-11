package handlers

import (
	"net/http"
)

// Ping check service HTTP handler
func (a *API) Ping(w http.ResponseWriter, r *http.Request) {
	if !a.usecase.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
