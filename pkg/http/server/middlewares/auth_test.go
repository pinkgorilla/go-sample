package middlewares_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinkgorilla/go-sample/pkg/auth"
	"github.com/pinkgorilla/go-sample/pkg/errors"
	"github.com/pinkgorilla/go-sample/pkg/http/server/middlewares"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	id := auth.FromContext(r.Context())
	if id == nil {
		err := errors.NewAuthError("id nil")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
	}
})

func Test_StaticAPIKeyValidatorMiddleware_ShouldOK(t *testing.T) {

	middleware := middlewares.AuthMiddleware(middlewares.StaticKeyAuthorizeFn("my-secret-api-key"))
	h := middleware(handler)
	server := httptest.NewServer(h)

	// createRequest: helper function to generate request with Bearer authorization
	createRequest := func(key string) *http.Request {
		req, err := http.NewRequest("GET", server.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))
		return req
	}

	// cases is map of request made and expected http status code from response
	cases := map[*http.Request]int{
		createRequest(""):                      http.StatusUnauthorized,
		createRequest("not-my-secret-api-key"): http.StatusUnauthorized,
		createRequest("my-secret-api-key"):     http.StatusOK,
	}

	for req, code := range cases {
		res, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != code {
			t.Fatalf("expected status %v, got %v", code, res.StatusCode)
		}
	}
}

func Test_StaticAPIKeyValidatorMiddleware_ShouldUnauthorized(t *testing.T) {
	middleware := middlewares.AuthMiddleware(middlewares.StaticKeyAuthorizeFn("my-secret-api-key"))
	h := middleware(handler)
	server := httptest.NewServer(h)

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	// req.Header.Set("Authorization", "Bearer my-secret-api-key")

	res, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %v, got %v", http.StatusUnauthorized, res.StatusCode)
	}
}
func Test_EmptyAPIKeyValidatorMiddleware_ShouldOK(t *testing.T) {
	middleware := middlewares.AuthMiddleware(middlewares.EmptyKeyAuthorizeFn)
	h := middleware(handler)
	server := httptest.NewServer(h)

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %v, got %v", http.StatusOK, res.StatusCode)
	}
}
