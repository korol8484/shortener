package handlers

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
	"github.com/korol8484/shortener/internal/app/user/util"
)

const (
	mimeJSON  = "application/json"
	mimePlain = "text/plain"
)

// Config Return HTTP domain append to short URL
type Config interface {
	GetBaseShortURL() string
}

// Store Repository Interface
type Store interface {
	Add(ctx context.Context, ent *domain.URL, user *domain.User) error
	Read(ctx context.Context, alias string) (*domain.URL, error)
	ReadByURL(ctx context.Context, URL string) (*domain.URL, error)
	AddBatch(ctx context.Context, batch domain.BatchURL, user *domain.User) error
	ReadUserURL(ctx context.Context, user *domain.User) (domain.BatchURL, error)
	BatchDelete(ctx context.Context, aliases []string, userID int64) error
	Close() error
}

// API api handler
type API struct {
	store Store
	cfg   Config
}

// NewAPI Factory
func NewAPI(store Store, cfg Config) *API {
	return &API{store: store, cfg: cfg}
}

// HandleShort Handler for one URL requested at plain text
// Response text/plain short URL
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

	userID, ok := util.ReadUserIDFromCtx(r.Context())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = a.store.Add(r.Context(), ent, &domain.User{ID: userID}); err != nil {
		if errors.Is(err, storage.ErrIssetURL) {
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

// HandleRedirect Handler plain text alias
// Response HTTP redirect to short URL
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

	if ent.Deleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	http.Redirect(w, r, ent.URL, http.StatusTemporaryRedirect)
}

func (a *API) shortURL(shortURL string) (*domain.URL, error) {
	parsedURL, err := url.Parse(shortURL)
	if err != nil {
		return nil, err
	}

	ent := &domain.URL{
		URL:   parsedURL.String(),
		Alias: GenAlias(6, shortURL),
	}

	return ent, nil
}

// GenAlias - Create alias length n as a hash of the string
func GenAlias(keyLen int, shortURL string) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	h := fnv.New64()
	h.Write([]byte(shortURL))

	r := rand.New(rand.NewSource(int64(h.Sum64())))

	keyMap := make([]byte, keyLen)
	for i := range keyMap {
		keyMap[i] = charset[r.Intn(len(charset))]
	}

	return string(keyMap)
}
