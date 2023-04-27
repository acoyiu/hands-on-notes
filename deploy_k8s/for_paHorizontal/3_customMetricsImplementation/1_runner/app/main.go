package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of the http requests received since the server started",
		},
		[]string{
			"namespace",
			"pod",
		},
	)
)

func init() {
	prometheus.MustRegister(HTTPRequests)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// code := 200
		switch path {
		case "/test":
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		case "/metrics":
			promhttp.Handler().ServeHTTP(w, r)
		default:
			w.WriteHeader(404)
			w.Write([]byte("Not Found"))
		}

		// Increace counter when request handler called
		HTTPRequests.
			With(prometheus.Labels{
				"namespace": os.Getenv("MY_NAMESPACE"),
				"pod":       os.Getenv("MY_POD_NAME"),
			}).
			Inc()
	})
	fmt.Println("Starting")
	http.ListenAndServe(":8080", nil)
}
