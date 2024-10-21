package handlers

import (
	"bytes"
	"context"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"github.com/korol8484/shortener/internal/app/user/storage"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_HandleShort(t *testing.T) {
	type want struct {
		code        int
		method      string
		contentType string
		body        string
	}
	tests := []struct {
		name string
		want want
	}{
		{name: "success_ya", want: want{
			method:      http.MethodPost,
			code:        201,
			contentType: "text/plain; charset=utf-8",
			body:        "http://www.ya.ru",
		}},
		{name: "not_post_request", want: want{
			method: http.MethodGet,
			code:   405,
		}},
		{name: "invalid_url_in_request", want: want{
			method: http.MethodPost,
			code:   400,
			body:   "http__://www.ya.ru",
		}},
	}

	router := chi.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()

	j := middleware.NewJwt(storage.NewMemoryStore(), zap.L(), "123")
	router.Use(j.HandlerSet())

	store := memory.NewMemStore()

	defer func(store Store) {
		_ = store.Close()
	}(store)

	api := NewAPI(store, &config.App{BaseShortURL: srv.URL})
	router.Post("/", api.HandleShort)

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.want.method, srv.URL+"/", strings.NewReader(test.want.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "text/plain; charset=utf-8")

			defer req.Body.Close()

			res, err := client.Do(req)
			require.NoError(t, err)

			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestAPI_HandleRedirect(t *testing.T) {
	router := chi.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()

	store := memory.NewMemStore()

	defer func(store Store) {
		_ = store.Close()
	}(store)

	api := NewAPI(store, &config.App{BaseShortURL: srv.URL})
	router.Get("/{id}", api.HandleRedirect)

	err := api.store.Add(context.Background(), &domain.URL{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	}, &domain.User{ID: 1})
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	type want struct {
		code        int
		method      string
		expectedURL string
		alias       string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "success",
			want: want{
				code:        http.StatusTemporaryRedirect,
				method:      http.MethodGet,
				expectedURL: "http://www.ya.ru",
				alias:       "7A2S4z",
			},
		},
		{
			name: "error bad method",
			want: want{
				code:        http.StatusMethodNotAllowed,
				method:      http.MethodPost,
				expectedURL: "",
				alias:       "7A2S4z",
			},
		},
		{
			name: "alias not found",
			want: want{
				code:        http.StatusBadRequest,
				method:      http.MethodGet,
				expectedURL: "",
				alias:       "111111",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, err := url.JoinPath(srv.URL, test.want.alias)
			require.NoError(t, err)

			req, err := http.NewRequest(test.want.method, path, nil)
			require.NoError(t, err)

			res, err := client.Do(req)
			require.NoError(t, err)

			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.expectedURL, res.Header.Get("Location"))
		})
	}
}

func TestAPI_ShortenJson(t *testing.T) {
	router := chi.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()

	j := middleware.NewJwt(storage.NewMemoryStore(), zap.L(), "123")
	router.Use(j.HandlerSet())

	store := memory.NewMemStore()

	defer func(store Store) {
		_ = store.Close()
	}(store)

	api := NewAPI(store, &config.App{BaseShortURL: srv.URL})
	router.Post("/json", api.ShortenJSON)

	client := &http.Client{}

	type want struct {
		code        int
		method      string
		contentType string
		body        string
	}
	tests := []struct {
		name string
		want want
	}{
		{name: "success_ya", want: want{
			method:      http.MethodPost,
			code:        201,
			contentType: "application/json",
			body:        "{\"url\": \"https://practicum.yandex.ru\"}",
		}},
		{name: "not_post_request", want: want{
			method: http.MethodGet,
			code:   405,
		}},
		{name: "invalid_url_in_request", want: want{
			method: http.MethodPost,
			code:   400,
			body:   "{\"url\": \"https__://practicum.yandex.ru\"}",
		}},
		{name: "success_ya_isset_url", want: want{
			method:      http.MethodPost,
			code:        409,
			contentType: "application/json",
			body:        "{\"url\": \"https://practicum.yandex.ru\"}",
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.want.method, srv.URL+"/json", strings.NewReader(test.want.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			defer req.Body.Close()

			res, err := client.Do(req)
			require.NoError(t, err)

			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func BenchmarkGenAlias(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenAlias(5, "https://ya.ru")
	}
}

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
