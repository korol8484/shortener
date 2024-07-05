package handlers

import (
	"github.com/korol8484/shortener/internal/app"
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
			code:   400,
		}},
		{name: "invalid url", want: want{
			method: http.MethodPost,
			code:   400,
			body:   "http__://www.ya.ru",
		}},
	}

	api := NewAPI(storage.NewMemStore())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.want.method, "http://localhost:8080/", strings.NewReader(test.want.body))
			w := httptest.NewRecorder()
			api.HandleShort(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			_ = res.Body.Close()
		})
	}
}

func TestAPI_HandleRedirect(t *testing.T) {
	api := NewAPI(storage.NewMemStore())

	err := api.store.Add(&app.Entity{
		URL:   "http://www.ya.ru",
		Alias: "7A2S4z",
	})
	require.NoError(t, err)

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
				code:        http.StatusBadRequest,
				method:      http.MethodPost,
				expectedURL: "",
				alias:       "7A2S4z",
			},
		},
		{
			name: "alias not found",
			want: want{
				code:        http.StatusBadRequest,
				method:      http.MethodPost,
				expectedURL: "",
				alias:       "111111",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, err := url.JoinPath("http://localhost:8080", test.want.alias)
			require.NoError(t, err)

			request := httptest.NewRequest(test.want.method, path, nil)
			w := httptest.NewRecorder()
			api.HandleRedirect(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.expectedURL, res.Header.Get("Location"))

			_ = res.Body.Close()
		})
	}
}
