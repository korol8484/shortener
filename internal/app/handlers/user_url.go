package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/user/util"
)

type responseURL struct {
	URL   string `json:"short_url"`
	Alias string `json:"original_url"`
}

// UserURL Handler for list user shortened links
// Returns:
//
//	[{
//	    "short_url": "http://localhost:8080/ZyNJrg",
//		"original_url": "http://ya.ru"
//	}]
func (a *API) UserURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := util.ReadUserIDFromCtx(r.Context())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	batch, err := a.store.ReadUserURL(r.Context(), &domain.User{
		ID: userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(batch) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]*responseURL, 0, len(batch))
	for _, u := range batch {
		resp = append(resp, &responseURL{
			URL:   fmt.Sprintf("%s/%s", a.cfg.GetBaseShortURL(), u.Alias),
			Alias: u.URL,
		})
	}

	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", mimeJSON)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
