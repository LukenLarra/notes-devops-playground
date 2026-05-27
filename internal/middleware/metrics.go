package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/luken/notes-devops-playground/internal/metrics"
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.ActiveRequests.Inc()
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		metrics.ActiveRequests.Dec()

		metrics.RequestCount.WithLabelValues(
			r.Method, r.URL.Path, strconv.Itoa(rw.status),
		).Inc()

		metrics.RequestDuration.WithLabelValues(
			r.Method, r.URL.Path,
		).Observe(duration)
	})
}
