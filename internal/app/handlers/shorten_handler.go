package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
	"github.com/korol8484/shortener/internal/app/user/util"
)

type request struct {
	URL string `json:"url"`
}

type response struct {
	Result string `json:"result"`
}

func (a *API) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// по сути лишнее, закрывается в net/http
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	req := &request{}

	if err = json.Unmarshal(body, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ent, err := a.shortURL(req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userId := util.ReadUserIdFromCtx(r.Context())
	if err = a.store.Add(r.Context(), ent, &domain.User{ID: userId}); err != nil {
		if errors.Is(err, storage.ErrIssetURL) {
			res := &response{Result: fmt.Sprintf("%s/%s", a.cfg.GetBaseShortURL(), ent.Alias)}

			b, err := json.Marshal(res)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("content-type", mimeJSON)
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write(b)

			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := &response{Result: fmt.Sprintf("%s/%s", a.cfg.GetBaseShortURL(), ent.Alias)}

	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", mimeJSON)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(b)
}
