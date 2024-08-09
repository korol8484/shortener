package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/korol8484/shortener/internal/app/handlers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAPI_Ping(t *testing.T) {
	type want struct {
		code   int
		method string
	}

	router := chi.NewRouter()
	srv := httptest.NewServer(router)
	defer srv.Close()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockPingable(ctrl)
	router.Get("/ping", Ping(m))

	gomock.InOrder(
		m.EXPECT().Ping().Return(errors.New("error")),
		m.EXPECT().Ping().Return(nil),
	)

	tests := []struct {
		name string
		want want
	}{
		{
			name: "ping_error",
			want: struct {
				code   int
				method string
			}{code: http.StatusInternalServerError, method: http.MethodGet},
		},
		{
			name: "ping_ok",
			want: struct {
				code   int
				method string
			}{code: http.StatusOK, method: http.MethodGet},
		},
	}

	client := &http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, err := url.JoinPath(srv.URL, "/ping")
			require.NoError(t, err)

			req, err := http.NewRequest(test.want.method, path, nil)
			require.NoError(t, err)

			res, err := client.Do(req)
			require.NoError(t, err)

			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
