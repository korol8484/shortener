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

const (
	mimeJson  = "application/json"
	mimePlain = "text/plain"
)

type Config interface {
	GetBaseShortURL() string
}

type API struct {
	store storage.Store
	cfg   Config
}

func NewAPI(store storage.Store, cfg Config) *API {
	return &API{store: store, cfg: cfg}
}

func (a *API) HandleShort(w http.ResponseWriter, r *http.Request) {
	cT := filterContentType(r.Header.Get("Content-Type"))
	if cT != mimePlain {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// по сути лишнее, закрывается в net/http
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	ent, err := a.shortUrl(string(body))
	if err != nil {
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

func (a *API) shortUrl(shortUrl string) (*domain.URL, error) {
	parsedURL, err := url.Parse(shortUrl)
	if err != nil {
		return nil, err
	}

	ent := &domain.URL{
		URL:   parsedURL.String(),
		Alias: a.genAlias(6),
	}

	if err = a.store.Add(ent); err != nil {
		return nil, err
	}

	return ent, nil
}

func filterContentType(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}
