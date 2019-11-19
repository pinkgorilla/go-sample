package middlewares

import (
	"net/http"
	"strings"

	"github.com/pinkgorilla/go-sample/pkg/auth"
	"github.com/pinkgorilla/go-sample/pkg/errors"
)

// StaticAPIKeyValidator is Validator for bearer static api key ...
type StaticAPIKeyValidator struct {
	StaticKey string
}

// NewStaticAPIKeyValidator ...
func NewStaticAPIKeyValidator(staticKey string) Validator {
	return &StaticAPIKeyValidator{
		StaticKey: staticKey,
	}
}

// Validate validates key against a store
func (m *StaticAPIKeyValidator) Validate(w http.ResponseWriter, r *http.Request) (*auth.Identity, error) {
	key := m.getAuthKey(r)
	if key == "" || key != m.StaticKey {
		return nil, errors.NewAuthError("unauthorized")
	}
	return &auth.Identity{
		ID:   "shared",
		Type: "static-key",
		Name: "static",
	}, nil
}

func (m *StaticAPIKeyValidator) getAuthKey(r *http.Request) string {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer")

	if len(splitToken) < 2 {
		return ""
	}
	token = strings.Trim(splitToken[1], " ")
	return token
}
