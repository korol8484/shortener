package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	store storage.Store
}

func NewAPI(store storage.Store) *API {
	return &API{store: store}
}

func (a *API) HandleShort(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedURL, err := url.Parse(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ent := &domain.Entity{
		URL:   parsedURL.String(),
		Alias: a.genAlias(6),
	}

	err = a.store.Add(ent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://%s/%s", r.Host, ent.Alias)))
}

func (a *API) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")

	if alias == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ent, err := a.store.Read(alias)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, ent.URL, http.StatusTemporaryRedirect)
}

func (a *API) genAlias(keyLen int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	keyMap := make([]byte, keyLen)
	for i := range keyMap {
		keyMap[i] = charset[r.Intn(len(charset))]
	}

	return string(keyMap)
}
