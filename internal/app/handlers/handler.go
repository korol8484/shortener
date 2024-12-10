package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"

	"github.com/korol8484/shortener/internal/app/storage"
	"github.com/korol8484/shortener/internal/app/usecase"
	"github.com/korol8484/shortener/internal/app/user/util"
)

const (
	mimeJSON  = "application/json"
	mimePlain = "text/plain"
)

// API api handler
type API struct {
	usecase *usecase.Usecase
}

// NewAPI Factory
func NewAPI(usecase *usecase.Usecase) *API {
	return &API{usecase: usecase}
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

	userID, ok := util.ReadUserIDFromCtx(r.Context())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ent, err := a.usecase.CreateURL(r.Context(), string(body), userID)

	if err != nil {
		if errors.Is(err, storage.ErrIssetURL) {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(a.usecase.FormatAlias(ent)))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(a.usecase.FormatAlias(ent)))
}

// HandleRedirect Handler plain text alias
// Response HTTP redirect to short URL
func (a *API) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")

	if alias == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ent, err := a.usecase.LoadByAlias(r.Context(), alias)
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
