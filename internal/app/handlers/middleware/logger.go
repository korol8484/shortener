package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type writer struct {
	http.ResponseWriter
	code        int
	bytes       int
	wroteHeader bool
}

// LoggRequest - Сведения о запросах должны содержать:
// URI, метод запроса и время, затраченное на его выполнение.
func LoggRequest(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			defer func() {
				duration := time.Since(start)
				logger.Info(
					"request",
					zap.String("RequestURI", r.RequestURI),
					zap.String("Method", r.Method),
					zap.Duration("ElapsedTime", duration),
				)
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// LoggResponse - Сведения об ответах должны содержать:
// код статуса и размер содержимого ответа.
func LoggResponse(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wr := &writer{ResponseWriter: w}

			defer func() {
				logger.Info(
					"response",
					zap.Int("StatusCode", wr.code),
					zap.Int("BytesWritten", wr.bytes),
				)
			}()

			next.ServeHTTP(wr, r)
		})
	}
}

func (w *writer) WriteHeader(code int) {
	if !w.wroteHeader {
		w.code = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *writer) Write(buf []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(buf)
	w.bytes += n

	return n, err
}
