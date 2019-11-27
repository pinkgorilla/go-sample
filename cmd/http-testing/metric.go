package main

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func httpRequestsTotalCounterVec(method string, path string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
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
			Name:        "http_request_duration_seconds",
			Help:        "A histogram of latencies for requests.",
			Buckets:     []float64{.25, .5, 1, 2.5, 5, 10},
			ConstLabels: prometheus.Labels{"path": path, "method": method},
		},
		[]string{"handler"},
	)

	prometheus.MustRegister(duration)
	return duration
}

// Instrument is a middleware to instrument prometheus RED metrics
func Instrument(method string, path string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerDuration(
			httpRequestDuration(method, path).MustCurryWith(prometheus.Labels{"handler": path}),
			promhttp.InstrumentHandlerCounter(
				httpRequestsTotalCounterVec(method, path),
				next,
			),
		)
	}
}

type Metric struct {
	once sync.Once
	h    http.Handler
	c    *chi.Mux
}

func NewMetric(m *chi.Mux) *Metric {
	return &Metric{
		c: m,
	}
}

func (m *Metric) Measure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.once.Do(func() {
			m.h = Instrument(r.Method, r.URL.Path)(next)
		})
		m.h.ServeHTTP(w, r)
	})
}

func (m *Metric) IstrumentedMethod(method, pattern string, handler http.Handler) {
	m.c.With(Instrument(method, pattern)).Method(method, pattern, handler)
}

func (m *Metric) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.c.ServeHTTP(w, r)
}
