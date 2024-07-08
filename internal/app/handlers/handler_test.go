package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAPI_HandleShort(t *testing.T) {
	type want struct {
		code        int
		method      string
		response    string
		contentType string
		body        string
	}
	tests := []struct {
		name string
		want want
	}{
		{name: "success ya", want: want{
			method:      http.MethodPost,
			code:        201,
			contentType: "text/plain; charset=utf-8",
			body:        "http://www.ya.ru",
		}},
		{name: "not post", want: want{
			method: http.MethodGet,
			code:   405,
		}},
		{name: "invalid url", want: want{
			method: http.MethodPost,
			code:   400,
			body:   "http__://www.ya.ru",
		}},
	}

	api := NewAPI(storage.NewMemStore())

	router := chi.NewRouter()
	router.Post("/", api.HandleShort)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.want.method, srv.URL+"/", strings.NewReader(test.want.body))
			require.NoError(t, err)

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
	api := NewAPI(storage.NewMemStore())

	err := api.store.Add(&domain.Entity{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	})
	require.NoError(t, err)

	router := chi.NewRouter()
	router.Get("/{id}", api.HandleRedirect)

	srv := httptest.NewServer(router)
	defer srv.Close()

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
