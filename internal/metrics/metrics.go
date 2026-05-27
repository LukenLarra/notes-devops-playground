package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "notes_api_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "notes_api_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	ActiveRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "notes_api_active_requests",
		Help: "Number of active HTTP requests",
	})
)
