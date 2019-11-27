package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pinkgorilla/go-sample/pkg/auth"
	"github.com/pinkgorilla/go-sample/pkg/errors"
)

type AuthorizeFn func(w http.ResponseWriter, r *http.Request) (*auth.Identity, error)

// AuthMiddleware ...
func AuthMiddleware(fn AuthorizeFn) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			id, err := fn(w, r)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			}
			ctx := auth.ToContext(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// StaticKeyAuthorizeFn returns AuthorizeFn which validates key against Authorization bearer token
func StaticKeyAuthorizeFn(key string) AuthorizeFn {
	return AuthorizeFn(func(w http.ResponseWriter, r *http.Request) (*auth.Identity, error) {
		token := r.Header.Get("Authorization")
		splitToken := strings.Split(token, "Bearer")

		if len(splitToken) < 2 {
			return nil, errors.NewAuthError("invalid auth token")
		}
		token = strings.Trim(splitToken[1], " ")

		if token == "" || token != key {
			return nil, errors.NewAuthError("unauthorized")
		}

		return &auth.Identity{
			ID:   0,
			Type: "static-key",
			Name: "static-key",
		}, nil
	})
}

// EmptyKeyAuthorizeFn is AuthorizeFn which ignores authorization mechanism
var EmptyKeyAuthorizeFn = AuthorizeFn(func(w http.ResponseWriter, r *http.Request) (*auth.Identity, error) {
	return &auth.Identity{
		ID:   0,
		Type: "empty-key",
		Name: "anonymous",
	}, nil
})
