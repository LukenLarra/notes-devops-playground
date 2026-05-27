package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rw, r)

			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration", time.Since(start).String(),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			}
			if reqID := middleware.GetReqID(r.Context()); reqID != "" {
				attrs = append(attrs, "request_id", reqID)
			}

			logger.Info("request", attrs...)
		})
	}
}
