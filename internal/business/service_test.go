package business_test

import (
	"testing"

	"github.com/pinkgorilla/go-sample/internal/business"
	"github.com/pinkgorilla/go-sample/internal/external/auth"
	"github.com/pinkgorilla/go-sample/internal/external/auth/dummy"
)

var dummyAuthProvider auth.Provider
var authService *auth.Service
var businessService *business.Service

func init() {
	dummyAuthProvider = dummy.NewProvider()
	authService, err := auth.NewService(dummyAuthProvider)
	if err != nil {
		panic(err)
	}
	businessService = business.NewService(authService)
}

func TestGetToken(t *testing.T) {
	token := businessService.GetToken("user", "secret")
	if token == "" {
		t.Error("expected token")
	}
}

func TestGetToken2(t *testing.T) {
	token := businessService.GetToken("user", "secret")
	if token == "" {
		t.Error("expected token")
	}
}
