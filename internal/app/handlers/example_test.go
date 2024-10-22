package handlers

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/korol8484/shortener/internal/app/user/storage"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

// Example - Пример использования сервиса
func Example() {
	// Create chi router
	router := chi.NewRouter()

	srv := httptest.NewServer(router)
	defer srv.Close()

	// Create JWT midelware for auth user
	j := middleware.NewJwt(storage.NewMemoryStore(), zap.L(), "123")
	router.Use(j.HandlerSet())

	// Create short URL storage
	store := memory.NewMemStore()

	defer func(store Store) {
		_ = store.Close()
	}(store)

	// Create API Handlers
	api := NewAPI(store, &config.App{BaseShortURL: srv.URL})
	apiDelete, err := NewDelete(store, zap.L())
	if err != nil {
		log.Fatal(err)
	}

	defer apiDelete.Close()

	// Register handlers in router
	router.Post("/", api.HandleShort)
	router.Get("/{id}", api.HandleRedirect)
	router.Post("/json", api.ShortenJSON)
	router.Delete("/batch", apiDelete.BatchDelete)
	router.Post("/batch", api.ShortenBatch)

	// Example clients request
	if code, _ := testRequest(srv, "POST", "/", nil, bytes.NewReader([]byte("http://ya.ru"))); code != 201 {
		log.Fatal("error create short url")
	}
}

func testRequest(ts *httptest.Server, method, path string, header http.Header, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		log.Fatal(err)
		return 0, ""
	}

	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return 0, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return 0, ""
	}
	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}
