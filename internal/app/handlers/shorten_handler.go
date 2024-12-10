package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/korol8484/shortener/internal/app/storage"
	"github.com/korol8484/shortener/internal/app/user/util"
)

type request struct {
	URL string `json:"url"`
}

type response struct {
	Result string `json:"result"`
}

// ShortenJSON Handler for json shortened link
// Accepts input json:
//
//	{
//	    "url": "http://www.ya.ru"
//	}
//
// Returns:
//
//	{
//	    "result": "http://localhost:8080/ZyNJrg"
//	}
func (a *API) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := &request{}
	if err = json.Unmarshal(body, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, ok := util.ReadUserIDFromCtx(r.Context())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ent, err := a.usecase.CreateURL(r.Context(), req.URL, userID)
	if err != nil {
		if errors.Is(err, storage.ErrIssetURL) {
			res := &response{Result: a.usecase.FormatAlias(ent)}

			var b []byte
			b, err = json.Marshal(res)
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

	res := &response{Result: a.usecase.FormatAlias(ent)}
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", mimeJSON)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(b)
}
