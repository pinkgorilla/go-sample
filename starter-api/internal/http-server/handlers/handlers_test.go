package handlers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pinkgorilla/go-sample/starter-api/internal/http-server/handlers"
)

var fn = func(w http.ResponseWriter, r *http.Request) error {
	time.Sleep(1 * time.Second)
	log.Println(r.Context().Err())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test:ok"))
	return nil
}

func Test_Handler_Request(t *testing.T) {
	executed := false
	server := httptest.NewServer(handlers.Handler(func(w http.ResponseWriter, r *http.Request) error {
		time.Sleep(500 * time.Millisecond)
		executed = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test:ok"))
		return nil
	}))

	c := server.Client()
	c.Timeout = 1 * time.Second
	res, err := c.Get(server.URL)

	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %v, got %v", http.StatusOK, res.StatusCode)
	}
	if !executed {
		t.Fatal("expected executed")
	}
}

func Test_Handler_RequestCancelled(t *testing.T) {
	executed := false
	server := httptest.NewServer(handlers.Handler(func(w http.ResponseWriter, r *http.Request) error {
		time.Sleep(500 * time.Millisecond)
		executed = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test:ok"))
		return nil
	}))

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	tr := &http.Transport{}
	c := server.Client()
	c.Transport = tr

	c.Do(req)
	tr.CancelRequest(req)
	if executed {
		t.Fatal("expected !executed")
	}
}
