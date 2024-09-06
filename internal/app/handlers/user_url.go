package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/user/util"
	"net/http"
)

type responseURL struct {
	URL   string `json:"short_url"`
	Alias string `json:"original_url"`
}

func (a *API) UserURL(w http.ResponseWriter, r *http.Request) {
	userID := util.ReadUserIDFromCtx(r.Context())

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
