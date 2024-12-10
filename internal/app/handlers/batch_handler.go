package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/user/util"
)

type batchRequestItem struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

type batchResponseItem struct {
	ID  string `json:"correlation_id"`
	URL string `json:"short_url"`
}

type batchRequest []batchRequestItem
type batchResponse []batchResponseItem

// ShortenBatch Handler for a collection of shortened links
// Accepts input json:
//
//	[{
//	    "correlation_id": "id",
//	    "original_url": "http://www.ya.ru"
//	}]
//
// Returns:
//
//	[{
//	    "correlation_id": "id",
//	    "short_url": "http://localhost:8080/ZyNJrg"
//	}]
func (a *API) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req batchRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(req) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	batchD := make(domain.BatchURL, 0, len(req))
	batchR := make(batchResponse, 0, len(req))

	for _, v := range req {
		var ent *domain.URL

		ent, err = a.usecase.GenerateURL(v.URL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		batchD = append(batchD, ent)
		batchR = append(batchR, batchResponseItem{
			ID:  v.ID,
			URL: a.usecase.FormatAlias(ent),
		})
	}

	userID, ok := util.ReadUserIDFromCtx(r.Context())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = a.usecase.AddBatch(r.Context(), batchD, userID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(batchR)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", mimeJSON)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(b)
}
