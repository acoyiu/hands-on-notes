package ginprometheus

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "http_server",
		Name:        "http_request_duration_seconds",
		Help:        "Histogram of response latency (seconds) of http handlers.",
		ConstLabels: nil,
		Buckets:     nil,
	}, []string{"method", "path"})

	httpReqTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests made.",
	}, []string{"method", "path", "status"})
)

func init() {
	prometheus.MustRegister(httpReqDuration)
	prometheus.MustRegister(httpReqTotal)
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		httpReqDuration.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.String(),
		}).Observe(time.Since(start).Seconds())

		httpReqTotal.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.String(),
			"status": strconv.Itoa(c.Writer.Status()),
		}).Inc()
	}
}
