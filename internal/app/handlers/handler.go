package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
)

const (
	mimeJSON  = "application/json"
	mimePlain = "text/plain"
)

type Config interface {
	GetBaseShortURL() string
}

type Store interface {
	Add(ctx context.Context, ent *domain.URL) error
	Read(ctx context.Context, alias string) (*domain.URL, error)
	ReadByURL(ctx context.Context, URL string) (*domain.URL, error)
	AddBatch(ctx context.Context, batch domain.BatchURL) error
	Close() error
}

type API struct {
	store Store
	cfg   Config
}

func NewAPI(store Store, cfg Config) *API {
	return &API{store: store, cfg: cfg}
}

func (a *API) HandleShort(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// по сути лишнее, закрывается в net/http
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	ent, err := a.shortURL(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = a.store.Add(r.Context(), ent); err != nil {
		if errors.Is(err, storage.ErrIssetURL) {
			ent, err = a.store.ReadByURL(r.Context(), ent.URL)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("%s/%s", a.cfg.GetBaseShortURL(), ent.Alias)))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", a.cfg.GetBaseShortURL(), ent.Alias)))
}

func (a *API) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")

	if alias == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ent, err := a.store.Read(r.Context(), alias)
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

func (a *API) shortURL(shortURL string) (*domain.URL, error) {
	parsedURL, err := url.Parse(shortURL)
	if err != nil {
		return nil, err
	}

	ent := &domain.URL{
		URL:   parsedURL.String(),
		Alias: a.genAlias(6),
	}

	return ent, nil
}
