package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpReqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	httpReqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpReqCount, httpReqDuration)
}

// PrometheusMetrics instruments Gin requests (request count + latency).
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Avoid self-instrumenting the metrics endpoint.
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		httpReqCount.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()
		httpReqDuration.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Observe(time.Since(start).Seconds())
	}
}

// PrometheusMiddleware allows creating a handler-friendly metrics middleware.
// Kept for potential future use.
func PrometheusMiddleware(_ http.Handler) gin.HandlerFunc {
	return PrometheusMetrics()
}
