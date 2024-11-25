package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriter(t *testing.T) {
	w := &writer{
		ResponseWriter: httptest.NewRecorder(),
		code:           0,
		bytes:          0,
		wroteHeader:    false,
	}

	_, err := w.Write([]byte("1"))
	require.NoError(t, err)

	w.WriteHeader(1)
}

func TestLoggRequest(t *testing.T) {
	withLogger(t, zapcore.DebugLevel, nil, func(logger *zap.Logger, logs *observer.ObservedLogs, t *testing.T) {
		r := chi.NewRouter()
		r.Use(LoggRequest(logger), LoggResponse(logger))

		r.Get("/get", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("htmlstring"))
		})

		ts := httptest.NewServer(r)
		defer ts.Close()

		req, err := http.NewRequest("GET", ts.URL+"/get", nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		defer resp.Body.Close()

		assert.Equal(t, logs.Len(), 2)
	})
}

func withLogger(t *testing.T, e zapcore.LevelEnabler, opts []zap.Option, f func(*zap.Logger, *observer.ObservedLogs, *testing.T)) {
	fac, logs := observer.New(e)
	log := zap.New(fac, opts...)
	f(log, logs, t)
}
