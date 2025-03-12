package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of requests processed by the MyApp web server.",
		},
		[]string{"path", "method", "status"},
	)

	ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_errors_total",
			Help: "Total number of error requests processed by the MyApp web server.",
		},
		[]string{"path", "method", "status"},
	)
)

func PrometheusInit() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(ErrorCount)
}

func TrackMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// Remove IDs and query parameters from the path
		c.Next()
		status := c.Writer.Status()

		// To reduce cardinality
		if status == 404 {
			path = "/404"
		}
		RequestCount.WithLabelValues(path, c.Request.Method, http.StatusText(status)).Inc()
		if status >= 400 {
			ErrorCount.WithLabelValues(path, c.Request.Method, http.StatusText(status)).Inc()
		}

	}
}
