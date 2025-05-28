package httphandlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

const responseWriterKey = "responseWriter"

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requesDuration)
}

var requestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "count-requests",
		Subsystem: "http",
		Name:      "request_total",
		Help:      "Total number of HTTP requests",
	},
	[]string{"path", "method"},
)

var requesDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "duration",
		Name:      "http_request_duration_seconds",
		Help:      "Duration of HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"method", "path"},
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start).Seconds()
		requestCount.WithLabelValues(r.URL.Path, r.Method).Inc()
		requesDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
