package middlewares

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
 
var reg map[string]bool
var regOnce sync.Once

func httpRequestsTotalCounterVec(method string, path string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "sample",
			Name:        "http_requests_total",
			Help:        "A counter for requests to the wrapped handler.",
			ConstLabels: prometheus.Labels{"path": path, "method": method},
		},
		[]string{"code"},
	) 
	prometheus.MustRegister(counter)

	return counter
}

func httpRequestDuration(method string, path string) *prometheus.HistogramVec {
	// duration is partitioned by the HTTP method and handler. It uses custom
	// buckets based on the expected request duration.
	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   "sample",
			Name:        "http_requests_duration",
			Help:        "A histogram of latencies for requests.",
			Buckets:     []float64{.25, .5, 1, 2.5, 5, 10},
			ConstLabels: prometheus.Labels{"path": path, "method": method},
		},
		[]string{"code"},
	)
	prometheus.MustRegister(duration)

	return duration
}

// Metrics is a middleware to instrument prometheus RED metrics
func Metrics(method string, path string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerDuration(
			httpRequestDuration(method, path),
			next,
		)
	}
}

// Metrics2 is a middleware to instrument prometheus RED metrics
func Metrics2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		promhttp.InstrumentHandlerDuration(
			httpRequestDuration(method, path).MustCurryWith(prometheus.Labels{"handler": path}),
			promhttp.InstrumentHandlerCounter(
				httpRequestsTotalCounterVec(method, path),
				next,
			),
		)
	})
}
