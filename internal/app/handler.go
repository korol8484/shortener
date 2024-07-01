package app

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type API struct {
	store Store
}

func NewAPI(store Store) *API {
	return &API{store: store}
}

func (a *API) HandleShort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	ent := &Entity{
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
	if r.Method != http.MethodGet {
		// разрешаем только GET-запросы
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.RequestURI == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// panic for invalid regex
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	alias := re.ReplaceAllString(r.RequestURI, "")

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
