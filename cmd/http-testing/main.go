package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinkgorilla/go-sample/pkg/http/server/instrument"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := chi.NewRouter()
	h := http.HandlerFunc(HTTPHandlerString)
	hc := instrument.ChiHandlerCreator(r)
	// http.Handle("/", PreHandle(PostHandle(h)))
	// http.Handle("/m/{id}", m.Measure(h))
	// http.Handle("/metrics", promhttp.Handler())
	hc("GET", "/", PreHandle(PostHandle(h)))
	hc("GET", "/m/{id}", h)
	hc("GET", "/metrics", promhttp.Handler())

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8880", nil))
}

func PreHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("pre-handler")
		next.ServeHTTP(w, r)
	})
}

func PostHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer log.Println("post-handler")
		next.ServeHTTP(w, r)
	})
}
