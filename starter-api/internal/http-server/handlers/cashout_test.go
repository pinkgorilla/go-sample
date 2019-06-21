package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinkgorilla/go-sample/starter-api/internal/http-server/handlers"

	"github.com/pinkgorilla/go-sample/starter-api/internal/app"
	"github.com/pinkgorilla/go-sample/starter-api/internal/test"
)

func Test_CreateCashOutHandler(t *testing.T) {
	a := test.App
	mw := app.InjectorMiddleware(a)
	h := mw(handlers.CreateCashOutHandler())
	server := httptest.NewServer(h)
	res, err := server.Client().Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %v, got %v", http.StatusOK, res.StatusCode)
	}
}
