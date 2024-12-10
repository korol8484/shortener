package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/korol8484/shortener/internal/app/config"
	"github.com/korol8484/shortener/internal/app/handlers/middleware"
	"github.com/korol8484/shortener/internal/app/storage/memory"
	"github.com/korol8484/shortener/internal/app/usecase"
	"github.com/korol8484/shortener/internal/app/user/storage"
)

func TestAPI_ShortenBatch(t *testing.T) {
	router := chi.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()

	j := middleware.NewJwt(usecase.NewJwt(storage.NewMemoryStore(), zap.L(), "123"), zap.L())
	router.Use(j.HandlerSet())

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

	api := NewAPI(uCase)
	router.Post("/batch", api.ShortenBatch)

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
		{name: "success_batch", want: want{
			method:      http.MethodPost,
			code:        201,
			contentType: "application/json",
			body:        "[{\"correlation_id\":\"id\",\"original_url\":\"http://www.ya.ru\"}]",
		}},
		{name: "not_post_request", want: want{
			method: http.MethodGet,
			code:   405,
		}},
		{name: "invalid_url_in_request", want: want{
			method: http.MethodPost,
			code:   400,
			body:   "[{\"correlation_id\":\"id\",\"original_url\":\"http___://www.ya.ru\"}]",
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.want.method, srv.URL+"/batch", strings.NewReader(test.want.body))
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
