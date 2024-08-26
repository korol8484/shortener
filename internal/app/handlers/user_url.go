package handlers

import (
	"encoding/json"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/user/util"
	"net/http"
)

type responseURL struct {
	URL   string `json:"short_url"`
	Alias string `json:"original_url"`
}

func (a *API) UserUrl(w http.ResponseWriter, r *http.Request) {
	userId := util.ReadUserIdFromCtx(r.Context())

	batch, err := a.store.ReadUserUrl(r.Context(), &domain.User{
		ID: userId,
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
			URL:   u.URL,
			Alias: u.Alias,
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
