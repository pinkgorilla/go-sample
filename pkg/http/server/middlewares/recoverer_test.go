package middlewares_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinkgorilla/go-sample/pkg/http/server/middlewares"
)

func Test_Recoverer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()

	panickingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("OMG!! i'm panicking!!")
	})

	s := httptest.NewServer(middlewares.Recoverer(panickingHandler))
	c := s.Client()
	data := map[string]interface{}{"name": "john", "age": 26}
	bs, err := json.Marshal(data)

	req, err := http.NewRequest("POST", fmt.Sprint(s.URL, "?x=1&y=2"), bytes.NewBuffer(bs))
	if err != nil {
		t.Fatal(err)
	}

	res, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected %v, got %v", http.StatusInternalServerError, res.StatusCode)
	}
}
