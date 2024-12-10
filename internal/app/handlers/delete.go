package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/korol8484/shortener/internal/app/user/util"
)

// BatchDelete Handler for a collection of delete user shorten URLs
// Accepts input json:
// ["cbi7jn", "dyifOs"]
// Returns: Http status Accepted
func (a *API) BatchDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := util.ReadUserIDFromCtx(r.Context())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var aliases []string

	if err = json.Unmarshal(body, &aliases); err != nil {
		http.Error(w, "can't parse json", http.StatusBadRequest)
		return
	}

	a.usecase.AddToDelete(aliases, userID)

	w.WriteHeader(http.StatusAccepted)
}
