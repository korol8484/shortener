package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/korol8484/shortener/internal/app/domain"
	"github.com/korol8484/shortener/internal/app/user/storage"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSample(t *testing.T) {
	r := chi.NewRouter()

	j := NewJwt(storage.NewMemoryStore(), zap.L(), "123")
	r.Use(j.HandlerSet(), j.HandlerRead())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	if status, resp := testRequest(t, ts, "GET", "/", nil, nil); status != 401 && resp == "welcome" {
		t.Fatalf(resp)
	}

	h := http.Header{}
	h.Set("Authorization", "BEARER asdf")
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 400 {
		t.Fatalf(resp)
	}

	jw, err := j.buildJWTString(&domain.User{1})
	if err != nil {
		t.Fatal(err)
	}

	h.Set("Authorization", "BEARER "+jw)
	if status, resp := testRequest(t, ts, "GET", "/", h, nil); status != 200 {
		t.Fatalf(resp)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, header http.Header, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}
	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}
