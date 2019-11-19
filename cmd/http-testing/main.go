package main

import (
	"log"
	"net/http"
)

func main() {
	h := http.HandlerFunc(HTTPHandlerString)
	http.Handle("/", PreHandle(PostHandle(h)))
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
