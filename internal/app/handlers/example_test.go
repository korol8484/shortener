package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/korol8484/shortener/internal/app/usecase"
	"github.com/korol8484/shortener/internal/app/user/storage"
)

// Example - Пример использования сервиса
func Example() {
	// Create chi router
	router := chi.NewRouter()

	srv := httptest.NewServer(router)
	defer srv.Close()

	// Create JWT midelware for auth user
	j := middleware.NewJwt(usecase.NewJwt(storage.NewMemoryStore(), zap.L(), "123"), zap.L())
	router.Use(j.HandlerSet())

	// Create short URL storage
	store := memory.NewMemStore()

	defer func(store usecase.Store) {
		_ = store.Close()
	}(store)

	uCase := usecase.NewUsecase(
		&config.App{BaseShortURL: srv.URL},
		store,
		usecase.NewPingDummy(),
		zap.L(),
	)
	defer uCase.Close()

	// Create API Handlers
	api := NewAPI(uCase)

	// Register handlers in router
	router.Post("/", api.HandleShort)
	router.Get("/{id}", api.HandleRedirect)
	router.Post("/json", api.ShortenJSON)
	router.Delete("/batch", api.BatchDelete)
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
